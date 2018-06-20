[![GoDoc](https://godoc.org/lufia.org/pkg/japanese/zipcode?status.svg)](https://godoc.org/lufia.org/pkg/japanese/zipcode)

# KEN_ALL.CSVパーサ

郵便番号データを扱いやすく整形します。

## Installation

```
go get lufia.org/pkg/japanese/zipcode
```

## Features

* [x] 複数行にまたがる郵便番号を1行に結合
* [x] 英数字や記号をASCII文字に統一
* [x] *以下に掲載がない場合* や *その他* などの文字を削除
* [x] *町域（ほげ、ふが）* を *町域ほげ*、*町域ふが* に分割
* [x] *町域（１〜３、５番地）* を *町域1番地*、*町域2番地* などに分割
* [ ] 地割対応
