package zipcode

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type entryParser interface {
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
	incompleteEntry = errors.New("incomplete entry")
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
				c1 <- incompleteEntry
				return
			}
			if err, ok := v.(error); ok {
				c1 <- err
				return
			}
			entry1 := v.(*Entry)
			entry.Town = entry.Town.combine(entry1.Town)
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
		remapRangeVerb(&entry.Town)
		a1, err := textRule.Eval(entry.Town.Text)
		if err != nil {
			c1 <- err
			return
		}
		a2, err := rubyRule.Eval(entry.Town.Ruby)
		if err != nil {
			c1 <- err
			return
		}

		// Town.Textには複数書式を持つが、Town.Rubyには複数部分を省略しているケースがある。
		if len(a1) > 1 && len(a2) == 1 {
			a3 := make([]string, len(a1))
			for i := 0; i < len(a1); i++ {
				a3[i] = a2[0]
			}
			a2 = a3
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

var parserFilters = []entryParser{
	entryHandlerFunc(func(entry *Entry) *Entry {
		if entry.Town.Text == "以下に掲載がない場合" {
			entry.Notice = entry.Town.Text
			entry.Town.Text = ""
			entry.Town.Ruby = ""
		}
		return entry
	}),
	entryHandlerFunc(func(entry *Entry) *Entry {
		const (
			textSuffix = "（次のビルを除く）"
			rubySuffix = "(ﾂｷﾞﾉﾋﾞﾙｦﾉｿﾞｸ)"
		)
		if strings.HasSuffix(entry.Town.Text, textSuffix) {
			entry.Town.Text = entry.Town.Text[0 : len(entry.Town.Text)-len(textSuffix)]
			entry.Town.Ruby = entry.Town.Ruby[0 : len(entry.Town.Ruby)-len(rubySuffix)]
		}
		return entry
	}),
	entryHandlerFunc(func(entry *Entry) *Entry {
		const (
			garbText = "（高層棟）"
			garbRuby = "(ｺｳｿｳﾄｳ)"
		)
		entry.Town.Text = strings.Replace(entry.Town.Text, garbText, "", -1)
		entry.Town.Ruby = strings.Replace(entry.Town.Ruby, garbRuby, "", -1)
		return entry
	}),
	entryHandlerFunc(func(entry *Entry) *Entry {
		const (
			textSuffix = "の次に番地がくる場合"
			rubySuffix = "ﾉﾂｷﾞﾆﾊﾞﾝﾁｶﾞｸﾙﾊﾞｱｲ"
		)
		if strings.HasSuffix(entry.Town.Text, textSuffix) {
			entry.Notice = entry.Town.Text
			town := Name{
				Text: entry.Town.Text[0 : len(entry.Town.Text)-len(textSuffix)],
				Ruby: entry.Town.Ruby[0 : len(entry.Town.Ruby)-len(rubySuffix)],
			}
			if strings.HasSuffix(entry.Region.Text, town.Text) {
				entry.Town.Text = ""
				entry.Town.Ruby = ""
			} else {
				entry.Town.Text = town.Text
				entry.Town.Ruby = town.Ruby
			}
		}
		return entry
	}),
	entryHandlerFunc(func(entry *Entry) *Entry {
		const (
			textSuffix = "一円"
			rubySuffix = "ｲﾁｴﾝ"
		)
		if strings.HasSuffix(entry.Town.Text, textSuffix) {
			if entry.Town.Text != textSuffix {
				entry.Notice = entry.Town.Text
				entry.Town.Text = ""
				entry.Town.Ruby = ""
			}
		}
		return entry
	}),
	entryHandlerFunc(func(entry *Entry) *Entry {
		entry.Town.Text = normalizeText(entry.Town.Text)
		return entry
	}),
	entryCollectorFunc(func(entry *Entry) bool {
		open := strings.Count(entry.Town.Text, "(")
		close := strings.Count(entry.Town.Text, ")")
		return open == close
	}),
	entryExpanderFunc(func(entry *Entry) []*Entry {
		return []*Entry{entry}
	}),
}

type Parser struct {
	Error error
}

// Parseは郵便番号データを流すチャネルを返す。
// エラーが途中で発生した場合、チャネルはclosedになりparser.Errorにエラーをセットする。
func (parser *Parser) Parse(r io.Reader) <-chan *Entry {
	c1 := make(chan interface{})
	go readFromCSVLoop(r, c1)
	for _, parser := range parserFilters {
		c2 := make(chan interface{})
		go parser.Parse(c1, c2)
		c1 = c2
	}

	c := make(chan *Entry)
	go func() {
		defer close(c)
		for v := range c1 {
			if err, ok := v.(error); ok {
				parser.Error = err
				return
			}
			c <- v.(*Entry)
		}
	}()
	return c
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
