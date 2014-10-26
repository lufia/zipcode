package postal

import (
	"bytes"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	actuals := []string{
		`01101,"060  ","0600000","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｲｶﾆｹｲｻｲｶﾞﾅｲﾊﾞｱｲ","北海道","札幌市中央区","以下に掲載がない場合",0,0,0,0,0,0`,
	}
	expects := []*Entry{
		&Entry{
			Code:            "01101",
			OldZip:          "060  ",
			Zip:             "0600000",
			Pref:            Name{"北海道", "ﾎｯｶｲﾄﾞｳ"},
			Region:          Name{"札幌市中央区", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ"},
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

func TestParseSplitted(t *testing.T) {
	actuals := []string{
		`02206,"01855","0185501","ｱｵﾓﾘｹﾝ","ﾄﾜﾀﾞｼ","ｵｸｾ(ｱｵﾌﾞﾅ､ｺﾀﾀﾐｲｼ､ﾄﾜﾀﾞ､ﾄﾜﾀﾞｺﾊﾝｳﾀﾙﾍﾞ､ﾄﾜﾀﾞｺﾊﾝﾈﾉｸﾁ､","青森県","十和田市","奥瀬（青撫、小畳石、十和田、十和田湖畔宇樽部、十和田湖畔子ノ口、",1,1,0,0,0,0`,
		`02206,"01855","0185501","ｱｵﾓﾘｹﾝ","ﾄﾜﾀﾞｼ","ﾄﾜﾀﾞｺﾊﾝﾔｽﾐﾔ)","青森県","十和田市","十和田湖畔休屋）",1,1,0,0,0,0`,
	}
	expects := []*Entry{
		&Entry{
			Code:            "02206",
			OldZip:          "01855",
			Zip:             "0185501",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬青撫", "ｵｸｾｱｵﾌﾞﾅ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "02206",
			OldZip:          "01855",
			Zip:             "0185501",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬小畳石", "ｵｸｾｺﾀﾀﾐｲｼ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "02206",
			OldZip:          "01855",
			Zip:             "0185501",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬十和田", "ｵｸｾﾄﾜﾀﾞ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "02206",
			OldZip:          "01855",
			Zip:             "0185501",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬十和田湖畔宇樽部", "ｵｸｾﾄﾜﾀﾞｺﾊﾝｳﾀﾙﾍﾞ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "02206",
			OldZip:          "01855",
			Zip:             "0185501",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬十和田湖畔子ノ口", "ｵｸｾﾄﾜﾀﾞｺﾊﾝﾈﾉｸﾁ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "02206",
			OldZip:          "01855",
			Zip:             "0185501",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬十和田湖畔休屋", "ｵｸｾﾄﾜﾀﾞｺﾊﾝﾔｽﾐﾔ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
	}
	parseTest(t, actuals, expects, "\n")
}

func TestParseOthers(t *testing.T) {
	actuals := []string{
		`02206,"03403","0340301","ｱｵﾓﾘｹﾝ","ﾄﾜﾀﾞｼ","ｵｸｾ(ｿﾉﾀ)","青森県","十和田市","奥瀬（その他）",1,1,0,0,0,0`,
	}
	expects := []*Entry{
		&Entry{
			Code:            "02206",
			OldZip:          "03403",
			Zip:             "0340301",
			Pref:            Name{"青森県", "ｱｵﾓﾘｹﾝ"},
			Region:          Name{"十和田市", "ﾄﾜﾀﾞｼ"},
			Town:            Name{"奥瀬", "ｵｸｾ"},
			IsPartialTown:   true,
			IsLargeTown:     true,
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
		select {
		case err := <-ec:
			t.Errorf("Parse() = %v; Expect not error", err)
		case entry := <-c:
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
			if !entry.Pref.Equal(expect.Pref) {
				t.Errorf("Parse(): Pref = %q; Expect %q", entry.Pref, expect.Pref)
			}
			if !entry.Region.Equal(expect.Region) {
				t.Errorf("Parse(): Region = %q; Expect %q", entry.Region, expect.Region)
			}
			if !entry.Town.Equal(expect.Town) {
				t.Errorf("Parse(): Town = %q; Expect %q", entry.Town, expect.Town)
			}
			if entry.IsPartialTown != expect.IsPartialTown {
				t.Errorf("Parse(): IsPartialTown = %q; Expect %q", entry.IsPartialTown, expect.IsPartialTown)
			}
			if entry.IsLargeTown != expect.IsLargeTown {
				t.Errorf("Parse(): IsLargeTown = %q; Expect %q", entry.IsLargeTown, expect.IsLargeTown)
			}
			if entry.IsBlockedScheme != expect.IsBlockedScheme {
				t.Errorf("Parse(): IsBlockedScheme = %q; Expect %q", entry.IsBlockedScheme, expect.IsBlockedScheme)
			}
			if entry.IsOverlappedZip != expect.IsOverlappedZip {
				t.Errorf("Parse(): IsOverlappedZip = %q; Expect %q", entry.IsOverlappedZip, expect.IsOverlappedZip)
			}
			if entry.Status != expect.Status {
				t.Errorf("Parse(): Status = %q; Expect %q", entry.Status, expect.Status)
			}
			if entry.Reason != expect.Reason {
				t.Errorf("Parse(): Reason = %q; Expect %q", entry.Reason, expect.Reason)
			}
		}
	}
}
