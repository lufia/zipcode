package postal

import (
	"io"
)

type ArrayFile struct {
	next    int
	records [][]string
}

func NewArrayFile(records [][]string) *ArrayFile {
	return &ArrayFile{0, records}
}

func (f *ArrayFile) Read() (record []string, err error) {
	if len(f.records) >= f.next {
		return nil, io.EOF
	}
	record = f.records[f.next]
	f.next++
	return
}
