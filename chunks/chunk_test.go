package chunks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkWriteAfterCloseFails(t *testing.T) {
	assert := assert.New(t)
	input := "abc"
	w := NewChunkWriter()
	_, err := w.Write([]byte(input))
	assert.NoError(err)

	assert.NoError(w.Close())
	assert.Panics(func() { w.Write([]byte(input)) }, "Write() after Close() should barf!")
}

func TestChunkWriteAfterChunkFails(t *testing.T) {
	assert := assert.New(t)
	input := "abc"
	w := NewChunkWriter()
	_, err := w.Write([]byte(input))
	assert.NoError(err)

	_ = w.Chunk()
	assert.Panics(func() { w.Write([]byte(input)) }, "Write() after Chunk() should barf!")
}

func TestChunkChunkCloses(t *testing.T) {
	assert := assert.New(t)
	input := "abc"
	w := NewChunkWriter()
	_, err := w.Write([]byte(input))
	assert.NoError(err)

	w.Chunk()
	assert.Panics(func() { w.Write([]byte(input)) }, "Write() after Close() should barf!")
}
