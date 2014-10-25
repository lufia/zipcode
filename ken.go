// 郵便番号パーサ。
// 県=pref
// 市=city
// 区=wards
// 郡=district
// 町=town
// 村=village
// 丁目=block(city districts?)
// 番地=number
// 号=extension
// 大字=divided
// 小字=sub-divided
// 建物名=building
package postal

import (
	"encoding/csv"
	"io"
	"strconv"
)

// 郵便番号データのエントリを表す。
type Entry struct {
	// 全国地方公共団体コード。
	Code string

	// 旧郵便番号(5桁)。
	OldZip string

	// 郵便番号(7桁)。
	Zip string

	// 都道府県名。
	Pref Name

	// 市区町村名。
	Region Name

	// 町域名。
	Town Name

	// 町域が2つ以上の郵便番号を持つ。
	IsPartialTown bool

	// 小字ごとに番地が起番されている。
	IsLargeTown bool

	// 丁目を有する町域。
	IsBlockedScheme bool

	// 1つの郵便番号で2つ以上の町域をあらわす。
	IsOverlappedZip bool

	// 更新の有無。
	Status Status

	// 更新理由。
	Reason Reason
}

// ルビ付き名前を表す。
type Name struct {
	// 漢字表記の名前。
	Text string

	// カナ表記の名前。
	Ruby string
}

// 名前が同じものかどうかを返す。
func (name Name) Equal(name1 Name) bool {
	return name.Text == name1.Text && name.Ruby == name1.Ruby
}

// 更新の表示。
type Status int

const (
	// 変更なし。
	StatusNotModified Status = 0

	// 変更あり。
	StatusModified = 1

	// 廃止。
	StatusObsoleted = 2
)

// 更新理由を表す。
type Reason int

const (
	ReasonNotModified Reason = 0
)

var parserChain = []Parser{}

// パースした結果を流すチャネルを返す。
// ecからエラーを受信した後はどちらのチャネルからもデータは届かない。
// 処理が終わればどちらのチャネルもクローズされる。
func Parse(c chan<- *Entry, ec chan<- error, r io.Reader) {
	c1 := make(chan interface{})
	go readFromCSVLoop(c1, r)
	for _, parser := range parserChain {
		c2 := make(chan interface{})
		parser.Parse(c1, c2)
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
func readFromCSVLoop(c chan<- interface{}, r io.Reader) {
	fin := csv.NewReader(r)
	record, err := fin.Read()
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

func parseStatus(s string) (Status, error) {
	switch s {
	case "0":
		return StatusNotModified, nil
	case "1":
		return StatusModified, nil
	case "2":
		return StatusObsoleted, nil
	default:
		return 0, strconv.ErrSyntax
	}
}

func parseReason(s string) (Reason, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	switch s {
	case "0", "1", "2", "3", "4", "5", "6":
		return Reason(n), nil
	default:
		return 0, strconv.ErrSyntax
	}
}
