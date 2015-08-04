package rabinkarp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRabinKarp(t *testing.T) {
	assert := assert.New(t)

	// This is the example from Wikipedia: https://en.wikipedia.org/wiki/Rabin%E2%80%93Karp_algorithm#Hash_function_used
	h := NewRabinKarp(3)
	h.HashByte('a')
	assert.Equal(uint32('a'), h.Sum32())
	h.HashByte('b')
	assert.Equal(uint32(9895), h.Sum32())
	h.HashByte('r')
	assert.Equal(uint32(999509), h.Sum32())
	h.HashByte('a')
	assert.Equal(uint32(1011309), h.Sum32())
}
