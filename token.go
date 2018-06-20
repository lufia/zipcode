package zipcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// cmplxRuleはKEN_ALL.CSVで複数の要素をまとめた書式をあらわす。
// 結合はRange>AddrSep>Delimの順で強い。
type cmplxRule struct {
	// 範囲書式の開始文字。たとえば"（"など。
	TokenBegin rune
	// 複数書式の終了文字。たとえば"）"など。
	TokenEnd rune
	// 複数書式で要素を分割する文字。たとえば"、"など。
	Delim rune
	// 複数書式で範囲をあらわす文字。たとえば"〜"など。
	Range rune
	// 1234-5のように番地の階層を区切る文字。
	AddrSep rune

	// 複数書式で範囲をあらわす文字の内部表現。
	// カナ文字はRangeとAddrSepが同じ文字を使うため置き換える。
	To rune
}

// Eval は、sに（）があれば内側の文字列を展開して、一連の文字列を配列で返す。
//
//	"あああ（ほげ、ふが）" => ["あああほげ", "あああふが"]
//	"（1〜3、5丁目）" => ["1丁目", "2丁目", "3丁目"]
func (rule cmplxRule) Eval(s string) ([]string, error) {
	t := []rune(s)
	for i := 0; i < len(t); i++ {
		switch t[i] {
		case rule.TokenBegin:
			p, err := rule.Expr(t, i+1)
			if err != nil {
				return nil, err
			}
			expr := t[i+1 : p]
			tokens, err := rule.Tokens(expr)
			if err != nil {
				return nil, err
			}
			a := make([]string, len(tokens))
			for j, token := range tokens {
				// "その他"だけは特別扱い
				if token == "その他" || token == "ｿﾉﾀ" {
					token = ""
				}
				a[j] = string(t[0:i]) + token + string(t[p+1:])
			}
			return a, nil
		case rule.TokenEnd:
			return nil, fmt.Errorf("missing %q", rule.TokenBegin)
		}
	}
	return []string{s}, nil
}

// Expr returns a index of the TokenEnd. If not, returns an error.
func (rule cmplxRule) Expr(s []rune, off int) (int, error) {
	for i := off; i < len(s); i++ {
		switch s[i] {
		case rule.TokenBegin:
			return -1, fmt.Errorf("missing %q", rule.TokenEnd)
		case rule.TokenEnd:
			return i, nil
		}
	}
	return -1, fmt.Errorf("missing %q", rule.TokenEnd)
}

// Tokensはrule.Delimとrule.Rangeを評価して展開後の文字列配列を返す。
func (rule cmplxRule) Tokens(expr []rune) (tokens []string, err error) {
	stage1 := rule.Split(expr)
	for _, s := range stage1 {
		stage2, err := rule.Expand(s)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, stage2...)
	}
	return
}

// Split は、exprをDelimで区切って配列を返す。
//
//	"あ、い、う" => ["あ", "い", "う"]
//	"２０〜２１-４番地" => ["２０〜２１-４番地"]
//	"１８-４、２０-４〜５番地" => ["１８-４番地", "２０-４〜５番地"]
func (rule cmplxRule) Split(expr []rune) []string {
	var (
		s      string
		peak   string
		ext    bool
		fields []string
	)
	for len(expr) > 0 {
		expr, s, ext = getToken(expr, rule.Delim)
		if !ext {
			fields = append(fields, s)
			continue
		}

		var a []string
		a = append(a, s)
		for len(expr) > 0 {
			expr, s, ext = getToken(expr, rule.Delim)
			if !ext {
				peak = s
				break
			}
			a = append(a, s)
		}
		first := []rune(a[0])
		n := skip(first, func(c rune) bool {
			return !unicode.IsDigit(c)
		})
		prefix := string(first[:n])
		a[0] = string(first[n:])

		last := []rune(a[len(a)-1])
		n = skip(last, func(c rune) bool {
			return unicode.IsDigit(c) || c == rule.To || c == rule.AddrSep
		})
		suffix := string(last[n:])
		a[len(a)-1] = string(last[:n])

		for _, s := range a {
			fields = append(fields, fmt.Sprintf("%s%s%s", prefix, s, suffix))
		}
		if peak != "" {
			fields = append(fields, peak)
			peak = ""
		}
	}
	return fields
}

func getToken(expr []rune, sep rune) ([]rune, string, bool) {
	var (
		s   strings.Builder
		ext bool
	)
	for i, c := range expr {
		if c == sep {
			return expr[i+1:], s.String(), ext
		}
		if unicode.IsDigit(c) {
			ext = true
		}
		s.WriteRune(c)
	}
	return nil, s.String(), ext
}

func skip(s []rune, f func(c rune) bool) int {
	for i, c := range s {
		if !f(c) {
			return i
		}
	}
	return len(s)
}

// Expand はRangeを展開して複数の文字列を返す。
//
//	"２０〜２１-４番地" => ["２０-４番地", "２１-４番地"]
//	"２０-４〜５番地" => ["２０-４番地", "２０-５番地"]
func (rule cmplxRule) Expand(token string) (tokens []string, err error) {
	to := string(rule.To)
	sep := string(rule.AddrSep)
	var prefix, s1, s2, suffix string

	re := regexp.MustCompile(`(\d+)` + to + `(\d+)(` + sep + `(\d+))?`)
	m := re.FindStringSubmatchIndex(token)
	if m != nil {
		prefix = token[0:m[2]]
		s1 = token[m[2]:m[3]]
		s2 = token[m[4]:m[5]]
		suffix = token[m[5]:]
	} else {
		re = regexp.MustCompile(`\d+` + sep + `(\d+)` + to + `(\d+)`)
		m = re.FindStringSubmatchIndex(token)
		if m == nil {
			tokens = append(tokens, token)
			return
		}
		prefix = token[0:m[2]]
		s1 = token[m[2]:m[3]]
		s2 = token[m[4]:m[5]]
		suffix = token[m[5]:]
	}
	bp, err := strconv.Atoi(s1)
	if err != nil {
		return
	}
	ep, err := strconv.Atoi(s2)
	if err != nil {
		return
	}
	for i := bp; i <= ep; i++ {
		tokens = append(tokens, fmt.Sprintf("%s%d%s", prefix, i, suffix))
	}
	return
}

var (
	textRule = cmplxRule{
		TokenBegin: '(',
		TokenEnd:   ')',
		Delim:      '、',
		Range:      '〜',
		AddrSep:    '−',
		To:         '〜',
	}
	rubyRule = cmplxRule{
		TokenBegin: '(',
		TokenEnd:   ')',
		Delim:      '､',
		Range:      '-',
		AddrSep:    '-',
		To:         '~',
	}
)

type tokenizer struct {
	buf strings.Builder
	s   []rune
}

func (t *tokenizer) Advance(a ...rune) rune {
	m := make(map[rune]struct{})
	for _, c := range a {
		m[c] = struct{}{}
	}
	for i, c := range t.s {
		if _, ok := m[c]; ok {
			t.s = t.s[i:]
			return c
		}
		t.buf.WriteRune(c)
	}
	return utf8.RuneError
}

func (t *tokenizer) Next() {
	t.buf.WriteRune(t.s[0])
	t.s = t.s[1:]
}

func (t *tokenizer) Replace(c rune) {
	t.buf.WriteRune(c)
	t.s = t.s[1:]
}

func (t *tokenizer) String() string {
	for _, c := range t.s {
		t.buf.WriteRune(c)
	}
	t.s = nil
	return t.buf.String()
}

// remapRangeVerb はカナの範囲文字を他の記号と重複しない文字に置き換える。
func remapRangeVerb(name *Name) {
	text := tokenizer{s: []rune(name.Text)}
	ruby := tokenizer{s: []rune(name.Ruby)}
	if text.Advance(textRule.TokenBegin) == utf8.RuneError {
		return
	}
	text.Next()

	if ruby.Advance(rubyRule.TokenBegin) == utf8.RuneError {
		return
	}
	ruby.Next()

scan:
	for {
		c := text.Advance(textRule.Range, textRule.AddrSep, textRule.TokenEnd)
		text.Next()
		switch c {
		case utf8.RuneError:
			break scan
		case textRule.Range:
			ruby.Advance(rubyRule.Range)
			ruby.Replace(rubyRule.To)
		case textRule.AddrSep:
			ruby.Advance(rubyRule.AddrSep)
			ruby.Next()
		case textRule.TokenEnd:
			ruby.Advance(rubyRule.TokenEnd)
			ruby.Next()
			break scan
		}
	}
	name.Text = text.String()
	name.Ruby = ruby.String()
}
