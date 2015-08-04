package rabinkarp

import (
	"bytes"
	"encoding/binary"
	"math"
)

// const prime32 = 16777619
const base uint32 = 101

// RabinKarp implements the hash.Hash32 interface and also has a function to write a single byte.
type RabinKarp struct {
	state    uint32
	buf      []byte
	n        uint32
	bufpos   uint32
	overflow bool
}

func NewRabinKarp(n uint32) *RabinKarp {
	rv := new(RabinKarp)
	rv.n = n
	rv.buf = make([]byte, n)
	rv.Reset()
	return rv
}

// HashByte updates the hash with a single byte and returns the resulting sum
func (rk *RabinKarp) HashByte(b byte) uint32 {
	if rk.bufpos == rk.n {
		rk.overflow = true
		rk.bufpos = 0
	}

	state := rk.state

	if rk.overflow {
		// ASCII a = 97, b = 98, r = 114.
		// hash("abr") = (97 × 101^2) + (98 × 101^1) + (114 × 101^0) = 999,509
		//             base   old hash    old 'a'         new 'a'
		// hash("bra") = [101 × (999,509 - (97 × 101^2))] + (97 × 101^0) = 1,011,309
		state = base*(state-uint32(rk.buf[rk.bufpos])*uint32(math.Pow(float64(base), float64(rk.n-1)))) + uint32(b)
	} else {
		state = base*state + uint32(b)
	}

	rk.buf[rk.bufpos] = b
	rk.bufpos++

	rk.state = state
	return state
}

// Write updates the hash with the bytes from slice p
func (rk *RabinKarp) Write(p []byte) (int, error) {
	for _, b := range p {
		rk.HashByte(b)
	}

	return len(p), nil
}

// Sum appends the (little endian) checksum to b
func (rk *RabinKarp) Sum(b []byte) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, rk.state)
	hash := buf.Bytes()
	for _, hb := range hash {
		b = append(b, hb)
	}

	return b
}

func (rk *RabinKarp) Reset() {
	rk.state = 0
	rk.bufpos = 0
	rk.overflow = false
}

func (rk *RabinKarp) Size() int {
	return 4
}

func (rk *RabinKarp) BlockSize() int {
	return int(rk.n)
}

func (rk *RabinKarp) Sum32() uint32 {
	return rk.state
}
