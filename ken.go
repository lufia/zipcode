// 郵便番号パーサ。
package postal

// 郵便番号データのエントリを表す。
type Address struct {
	// 郵便番号(7桁)。
	Code string

	// 都道府県。
	Pref NameRuby

	// 市区町村。
	City NameRuby

	// 町域。
	SpecificPart NameRuby

	// 町域が2つ以上の郵便番号を持つ。
	A bool

	// 小字ごとに番地が起番されている。
	B bool

	// 1つの郵便番号で2つ以上の町域をあらわす。
	C bool

	// 更新あり。
	IsModified bool

	// 更新理由。
	Reason Reason
}

// aとbを比較して内容が同じならtrueを返す。
func (a *Address) Equal(b *Address) bool {
	return a == b
}

// ルビ付き名前を表す。
type NameRuby struct {
	Name string
	Ruby string
}

// 更新理由を表す。
type Reason int

// 削除されている場合はtrueを返す。
func (r Reason) IsDeleted() bool {
	return int(r) == 6
}

// 郵便番号表のレコードを読む。
type RecordReader interface {
	Read() (record []string, err error)
}

// パースした結果を流すチャネルを返す。
func Parse(r RecordReader) chan *Address {
	c := make(chan *Address)
	close(c)
	return c
}
