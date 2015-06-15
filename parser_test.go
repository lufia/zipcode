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

		`26104,"604  ","6040983","ｷｮｳﾄﾌ","ｷｮｳﾄｼﾅｶｷﾞｮｳｸ","ｻｻﾔﾁｮｳ","京都府","京都市中京区","笹屋町（麩屋町通竹屋町下る、麩屋町通夷川上る、竹屋町通麩屋町西入、竹屋",0,0,0,0,0,0`,
		`26104,"604  ","6040983","ｷｮｳﾄﾌ","ｷｮｳﾄｼﾅｶｷﾞｮｳｸ","ｻｻﾔﾁｮｳ","京都府","京都市中京区","町通麩屋町東入、竹屋町通御幸町西入、夷川通麩屋町西入、夷川通麩屋町東入）",0,0,0,0,0,0`,
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

		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町麩屋町通竹屋町下る", "ｻｻﾔﾁｮｳ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町麩屋町通夷川上る", "ｻｻﾔﾁｮｳ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町竹屋町通麩屋町西入", "ｻｻﾔﾁｮｳ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町竹屋町通麩屋町東入", "ｻｻﾔﾁｮｳ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町竹屋町通御幸町西入", "ｻｻﾔﾁｮｳ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町夷川通麩屋町西入", "ｻｻﾔﾁｮｳ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: false,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "26104",
			OldZip:          "604  ",
			Zip:             "6040983",
			Pref:            Name{"京都府", "ｷｮｳﾄﾌ"},
			Region:          Name{"京都市中京区", "ｷｮｳﾄｼﾅｶｷﾞｮｳｸ"},
			Town:            Name{"笹屋町夷川通麩屋町東入", "ｻｻﾔﾁｮｳ"},
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

func TestPrase(t *testing.T) {
	actuals := []string{
		`01101,"064  ","0640930","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ﾐﾅﾐ30ｼﾞｮｳﾆｼ(9-11ﾁｮｳﾒ)","北海道","札幌市中央区","南三十条西（９〜１１丁目）",0,0,1,0,0,0`,
	}
	expects := []*Entry{
		&Entry{
			Code:            "01101",
			OldZip:          "064  ",
			Zip:             "0640930",
			Pref:            Name{"北海道", "ﾎｯｶｲﾄﾞｳ"},
			Region:          Name{"札幌市中央区", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ"},
			Town:            Name{"南三十条西9丁目", "ﾐﾅﾐ30ｼﾞｮｳﾆｼ9ﾁｮｳﾒ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: true,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "01101",
			OldZip:          "064  ",
			Zip:             "0640930",
			Pref:            Name{"北海道", "ﾎｯｶｲﾄﾞｳ"},
			Region:          Name{"札幌市中央区", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ"},
			Town:            Name{"南三十条西10丁目", "ﾐﾅﾐ30ｼﾞｮｳﾆｼ10ﾁｮｳﾒ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: true,
			IsOverlappedZip: false,
			Status:          StatusNotModified,
			Reason:          ReasonNotModified,
		},
		&Entry{
			Code:            "01101",
			OldZip:          "064  ",
			Zip:             "0640930",
			Pref:            Name{"北海道", "ﾎｯｶｲﾄﾞｳ"},
			Region:          Name{"札幌市中央区", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ"},
			Town:            Name{"南三十条西11丁目", "ﾐﾅﾐ30ｼﾞｮｳﾆｼ11ﾁｮｳﾒ"},
			IsPartialTown:   false,
			IsLargeTown:     false,
			IsBlockedScheme: true,
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
				t.Errorf("Parse() = nil; Expect %v", expect)
				continue
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
				t.Errorf("Parse(): IsPartialTown = %t; Expect %t", entry.IsPartialTown, expect.IsPartialTown)
			}
			if entry.IsLargeTown != expect.IsLargeTown {
				t.Errorf("Parse(): IsLargeTown = %t; Expect %t", entry.IsLargeTown, expect.IsLargeTown)
			}
			if entry.IsBlockedScheme != expect.IsBlockedScheme {
				t.Errorf("Parse(): IsBlockedScheme = %t; Expect %t", entry.IsBlockedScheme, expect.IsBlockedScheme)
			}
			if entry.IsOverlappedZip != expect.IsOverlappedZip {
				t.Errorf("Parse(): IsOverlappedZip = %t; Expect %t", entry.IsOverlappedZip, expect.IsOverlappedZip)
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
