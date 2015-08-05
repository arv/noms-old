package adler32

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdler32(t *testing.T) {
	assert := assert.New(t)

	// This is the example from Wikipedia: https://en.wikipedia.org/wiki/Adler-32#Example
	h := NewAdler32(64)

	h.HashByte('W')
	assert.Equal(88, int(h.s1))
	assert.Equal(88, int(h.s2))
	assert.Equal(uint32(88<<16+88), h.Sum32())

	h.HashByte('i')
	assert.Equal(193, int(h.s1))
	assert.Equal(281, int(h.s2))
	assert.Equal(uint32(193<<16+281), h.Sum32())

	h.HashByte('k')
	assert.Equal(300, int(h.s1))
	assert.Equal(581, int(h.s2))
	assert.Equal(uint32(300<<16+581), h.Sum32())

	h.HashByte('i')
	assert.Equal(405, int(h.s1))
	assert.Equal(986, int(h.s2))
	assert.Equal(uint32(405<<16+986), h.Sum32())

	h.HashByte('p')
	assert.Equal(517, int(h.s1))
	assert.Equal(1503, int(h.s2))
	assert.Equal(uint32(517<<16+1503), h.Sum32())

	h.HashByte('e')
	assert.Equal(618, int(h.s1))
	assert.Equal(2121, int(h.s2))
	assert.Equal(uint32(618<<16+2121), h.Sum32())

	h.HashByte('d')
	assert.Equal(718, int(h.s1))
	assert.Equal(2839, int(h.s2))
	assert.Equal(uint32(718<<16+2839), h.Sum32())

	h.HashByte('i')
	assert.Equal(823, int(h.s1))
	assert.Equal(3662, int(h.s2))
	assert.Equal(uint32(823<<16+3662), h.Sum32())

	h.HashByte('a')
	assert.Equal(920, int(h.s1))
	assert.Equal(4582, int(h.s2))
	assert.Equal(uint32(920<<16+4582), h.Sum32())
}
