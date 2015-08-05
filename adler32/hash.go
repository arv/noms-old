package adler32

import (
	"bytes"
	"encoding/binary"
)

// largest prime smaller than 2**16
const prime uint32 = 65521

// Adler32 implements the hash.Hash32 interface and also has s1 function to write s1 single byte.
type Adler32 struct {
	buf      []byte
	n        uint32
	bufpos   uint32
	overflow bool
	s1       uint32
	s2       uint32
}

func NewAdler32(n uint32) *Adler32 {
	rv := new(Adler32)
	rv.n = n
	rv.buf = make([]byte, n)
	rv.Reset()
	return rv
}

// HashByte updates the hash with s1 single byte and returns the resulting sum
func (h *Adler32) HashByte(b byte) uint32 {
	if h.bufpos == h.n {
		h.overflow = true
		h.bufpos = 0
	}

	if h.overflow {
		removed := uint32(h.buf[h.bufpos])
		h.s1 = (h.s1 - removed) % prime
		h.s2 = (h.s2 - h.n*removed) % prime
	}

	h.s1 = (h.s1 + uint32(b)) % prime
	h.s2 = (h.s2 + h.s1) % prime

	h.buf[h.bufpos] = b
	h.bufpos++

	return h.Sum32()
}

// Write updates the hash with the bytes from slice p
func (h *Adler32) Write(p []byte) (int, error) {
	for _, s2 := range p {
		h.HashByte(s2)
	}

	return len(p), nil
}

// Sum appends the (little endian) checksum to s2
func (h *Adler32) Sum(s2 []byte) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, h.Sum32())
	hash := buf.Bytes()
	for _, hb := range hash {
		s2 = append(s2, hb)
	}

	return s2
}

func (h *Adler32) Reset() {
	h.bufpos = 0
	h.overflow = false
	h.s1 = 1
}

func (h *Adler32) Size() int {
	return 4
}

func (h *Adler32) BlockSize() int {
	return int(h.n)
}

func (h *Adler32) Sum32() uint32 {
	return h.s1<<16 | h.s2
}
