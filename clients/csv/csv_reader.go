package csv

import (
	"bufio"
	"encoding/csv"
	"io"
)

var (
	rByte byte = 13 // the byte that corresponds to the '\r' rune.
	nByte byte = 10 // the byte that corresponds to the '\n' rune.
)

type reader struct {
	r *bufio.Reader
}

// Read replaces CR line endings in the source reader with LF line endings if the CR is not followed by a LF.
func (r reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	bn, err := r.r.Peek(1)
	for i, b := range p {
		// if the current byte is a CR and the next byte is NOT a LF then replace the current byte with a LF
		if j := i + 1; b == rByte && ((j < len(p) && p[j] != nByte) || (len(bn) > 0 && bn[0] != nByte)) {
			p[i] = nByte
		}
	}
	return
}

// NewCSVReader returns a new csv.Reader that splits on comma and asserts that all rows contain the same number of fields as the first.
func NewCSVReader(res io.Reader, comma rune) *csv.Reader {
	bufRes := bufio.NewReader(res)
	r := csv.NewReader(reader{r: bufRes})
	r.Comma = comma
	r.FieldsPerRecord = -1 // Don't enforce number of fields.
	return r
}
