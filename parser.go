package postal

import (
	"fmt"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Parser interface {
	// Parseは、cから値を受信したあと必要な加工を行い、c1へ送信する。
	Parse(c <-chan interface{}, c1 chan<- interface{})
}

type entryHandlerFunc func(entry *Entry) *Entry

func (f entryHandlerFunc) Parse(c <-chan interface{}, c1 chan<- interface{}) {
	defer close(c1)
	for v := range c {
		if err, ok := v.(error); ok {
			c1 <- err
			break
		}
		c1 <- f(v.(*Entry))
	}
}

var (
	// 行連結が未解決のままとなった場合のエラー
	IncompleteEntry = errors.New("incomplete entry")
)

// entryが完結している場合はtrue、次のエントリと連結する場合はfalse
type entryCollectorFunc func(entry *Entry) bool

func (f entryCollectorFunc) Parse(c <-chan interface{}, c1 chan<- interface{}) {
	defer close(c1)
	for v := range c {
		if err, ok := v.(error); ok {
			c1 <- err
			return
		}
		entry := v.(*Entry)
		for !f(entry) {
			v, ok := <-c
			if !ok {
				c1 <- IncompleteEntry
				return
			}
			if err, ok := v.(error); ok {
				c1 <- err
				return
			}
			entry1 := v.(*Entry)
			entry.Town = entry.Town.Combine(entry1.Town)
		}
		c1 <- entry
	}
}

type entryExpanderFunc func(entry *Entry) []*Entry

func (f entryExpanderFunc) Parse(c <-chan interface{}, c1 chan<- interface{}) {
	defer close(c1)
	for v := range c {
		if err, ok := v.(error); ok {
			c1 <- err
			return
		}
		entry := v.(*Entry)
		a1, err := tokenDelim{'（', '）', '、'}.Expand(entry.Town.Text)
		if err != nil {
			c1 <- err
			return
		}
		a2, err := tokenDelim{'(', ')', '､'}.Expand(entry.Town.Ruby)
		if err != nil {
			c1 <- err
			return
		}
		if len(a1) != len(a2) {
			// TODO
		}
		for i, _ := range a1 {
			entry1 := new(Entry)
			*entry1 = *entry
			entry1.Town = Name{a1[i], a2[i]}
			c1 <- entry1
		}
	}
}

type tokenDelim struct {
	TokenBegin rune
	TokenEnd rune
	Delim rune
}

// （）の中の文字列を展開して配列で返す。
// "ほげ、ふが" => ["ほげ", "ふが"]
func (delim tokenDelim) Expand(s string) ([]string, error) {
	t := []rune(s)
	for i := 0; i < len(t); i++ {
		switch t[i] {
		case delim.TokenBegin:
			expr, err := delim.Expr(t[i+1:])
			if err != nil {
				return nil, err
			}
			tokens := strings.FieldsFunc(string(expr), func(c rune) bool {
				return c == delim.Delim
			})
			a := make([]string, len(tokens))
			for j, token := range tokens {
				// "その他"だけは特別扱い
				if token == "その他" || token == "ｿﾉﾀ" {
					token = ""
				}
				a[j] = string(t[0:i]) + token + string(t[i+len(expr)+2:])
			}
			return a, nil
		case delim.TokenEnd:
			return nil, fmt.Errorf("missing %q", delim.TokenBegin)
		}
	}
	return []string{s}, nil
}

func (delim tokenDelim) Expr(s []rune) ([]rune, error) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case delim.TokenBegin:
			return nil, fmt.Errorf("missing %q", delim.TokenEnd)
		case delim.TokenEnd:
			return s[0:i], nil
		}
	}
	return nil, fmt.Errorf("missing %q", delim.TokenEnd)
}

var parserChain = []Parser{
	entryHandlerFunc(func(entry *Entry) *Entry {
		if entry.Town.Text == "以下に掲載がない場合" {
			entry.Town.Text = ""
		}
		if entry.Town.Ruby == "ｲｶﾆｹｲｻｲｶﾞﾅｲﾊﾞｱｲ" {
			entry.Town.Ruby = ""
		}
		return entry
	}),
	entryCollectorFunc(func(entry *Entry) bool {
		open := strings.Count(entry.Town.Text, "（")
		close := strings.Count(entry.Town.Text, "）")
		return open == close
	}),
	entryExpanderFunc(func(entry *Entry) []*Entry {
		return []*Entry{entry}
	}),
}

// パースした結果を流すチャネルを返す。
// ecからエラーを受信した後はどちらのチャネルからもデータは届かない。
func Parse(c chan<- *Entry, ec chan<- error, r io.Reader) {
	defer close(c)
	c1 := make(chan interface{})
	go readFromCSVLoop(r, c1)
	for _, parser := range parserChain {
		c2 := make(chan interface{})
		go parser.Parse(c1, c2)
		c1 = c2
	}
	for v := range c1 {
		if err, ok := v.(error); ok {
			ec <- err
			return
		}
		c <- v.(*Entry)
	}
}

// rからCSVデータを読み、cにエントリを送信する。
// エラーが発生した場合はecへエラーを送信する。
func readFromCSVLoop(r io.Reader, c chan<- interface{}) {
	defer close(c)
	fin := csv.NewReader(r)
	for {
		record, err := fin.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			c <- err
			return
		}

		isPartialTown, err := strconv.ParseBool(record[9])
		if err != nil {
			c <- err
			return
		}
		isLargeTown, err := strconv.ParseBool(record[10])
		if err != nil {
			c <- err
			return
		}
		isBlockedScheme, err := strconv.ParseBool(record[11])
		if err != nil {
			c <- err
			return
		}
		isOverlappedZip, err := strconv.ParseBool(record[12])
		if err != nil {
			c <- err
			return
		}
		status, err := parseStatus(record[13])
		if err != nil {
			c <- err
			return
		}
		reason, err := parseReason(record[14])
		if err != nil {
			c <- err
			return
		}
		c <- &Entry{
			Code:            record[0],
			OldZip:          record[1],
			Zip:             record[2],
			Pref:            Name{record[6], record[3]},
			Region:          Name{record[7], record[4]},
			Town:            Name{record[8], record[5]},
			IsPartialTown:   isPartialTown,
			IsLargeTown:     isLargeTown,
			IsBlockedScheme: isBlockedScheme,
			IsOverlappedZip: isOverlappedZip,
			Status:          status,
			Reason:          reason,
		}
	}
}
