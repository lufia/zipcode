package postal

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseOneLine(t *testing.T) {
	actuals := []string{
		`01101,"060  ","0600000","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｲｶﾆｹｲｻｲｶﾞﾅｲﾊﾞｱｲ","北海道","札幌市中央区","以下に掲>載がない場合",0,0,0,0,0,0`,
	}
	expects := []*Entry{
		&Entry{
			Code:            "01101",
			OldZip:          "060  ",
			Zip:             "0600000",
			Pref:            Name{"北海道", "ﾎｯｶｲﾄﾞｳ"},
			Region:          Name{"幌市中央区", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ"},
			Town:            Name{"", ""},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
	}
	parseTest(t, actuals, expects, "\n")
}
func parseTest(t *testing.T, actuals []string, expects []*Entry, newline string) {
	c := make(chan *Entry)
	ec := make(chan error)
	s := strings.Join(actuals, newline)
	fin := bytes.NewBufferString(s)
	go Parse(c, ec, fin)
	for _, expect := range expects {
		entry := <-c
		if entry == nil {
			t.Errorf("Parse() = nil; Expect not nil")
			return
		}
		if entry.Code != expect.Code {
			t.Errorf("Parse(): Code = %q; Expect %q", entry.Code, expect.Code)
		}
		if entry.OldZip != expect.OldZip {
			t.Errorf("Parse(): OldZip = %q; Expect %q", entry.OldZip, expect.OldZip)
		}
		if entry.Zip != expect.Zip {
			t.Errorf("Parse(): Zip = %q; Expect %q", entry.Zip, expect.Zip)
		}
		if entry.Pref != expect.Pref {
			t.Errorf("Parse(): Pref = %q; Expect %q", entry.Pref, expect.Pref)
		}
	}
}
