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
// 参考: http://en.wikipedia.org/wiki/Japanese_addressing_system
package zipcode // import "lufia.org/pkg/japanese/zipcode"

import (
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

	//has_chome

	// 1つの郵便番号で2つ以上の町域をあらわす。
	IsOverlappedZip bool

	// 更新の有無。
	Status Status

	// 更新理由。
	Reason Reason

	// 備考。このフィールドはKEN_ALL.CSVには存在しない。
	Notice string
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

// Combineはname1の内容をnameの後に追加する。
// 追加された状態のNameを返す。
func (name Name) Combine(name1 Name) Name {
	var ruby string
	if name.Ruby == name1.Ruby {
		ruby = name.Ruby
	} else {
		ruby = name.Ruby + name1.Ruby
	}
	return Name{
		Text: name.Text + name1.Text,
		Ruby: ruby,
	}
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
