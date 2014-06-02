package postal

import (
	"testing"
)

// Reasonの全メソッド結果がすべて正しいことを確認。
func TestReason(t *testing.T) {
	tab := []struct {
		r Reason
		d bool
	}{
		{Reason(0), false},
		{Reason(1), false},
		{Reason(2), false},
		{Reason(3), false},
		{Reason(4), false},
		{Reason(5), false},
		{Reason(6), true},
	}
	for _, item := range tab {
		if d := item.r.IsDeleted(); d != item.d {
			t.Errorf("IsDeleted() = %v, but expected %v\n", d, item.d)
		}
	}
}

// ファイルが空の場合は住所リストも空になる。
func TestEmpty(t *testing.T) {
	source := [][]string{}
	expect := []*Address{}
	testParser(t, source, expect)
}

// テスト用の便利メソッド。
// 入力データをパースした結果と期待値を比較してエラー報告する。
func testParser(t *testing.T, records [][]string, a []*Address) {
	result := make([]*Address, 0, len(a))
	c := Parse(NewArrayFile(records))
	for addr := range c {
		result = append(result, addr)
	}
	max := len(result)
	if len(result) != len(a) {
		t.Errorf("len(result) = %d; expect %d\n", len(result), len(a))
		if len(result) < len(a) {
			max = len(a)
		}
	}
	for i := 0; i < max; i++ {
		if i >= len(result) {
			t.Errorf("len(result) too short: result[%d] = %v\n", i, *result[i])
		} else if i >= len(a) {
			t.Errorf("len(result) too long: result[%d] = %v\n", i, *result[i])
		} else if !result[i].Equal(a[i]) {
			t.Errorf("result[%d] = %v expect %v\n", result[i], a[i])
		}
	}
}
