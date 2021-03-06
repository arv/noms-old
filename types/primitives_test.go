package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimitives(t *testing.T) {
	data := []Value{
		Bool(true), Bool(false),
		Int16(0), Int16(-1),
		Int32(0), Int32(-1),
		Int64(0), Int64(-1),
		Uint16(0), Uint16(1),
		Uint32(0), Uint32(1),
		Uint64(0), Uint64(1),
		Float32(0.0), Float32(0.1),
		Float64(0.0), Float64(0.1),
	}

	for i := range data {
		for j := range data {
			if i == j {
				assert.True(t, data[i].Equals(data[j]), "Expected value to equal self at index %d", i)
			} else {
				assert.False(t, data[i].Equals(data[j]), "Expected values at indices %d and %d to not equal", i, j)
			}
		}
	}
}

func TestPrimitivesType(t *testing.T) {
	data := []struct {
		v Value
		k NomsKind
	}{
		{Bool(false), BoolKind},
		{Int8(0), Int8Kind},
		{Int16(0), Int16Kind},
		{Int32(0), Int32Kind},
		{Int64(0), Int64Kind},
		{Float32(0), Float32Kind},
		{Float64(0), Float64Kind},
		{Uint8(0), Uint8Kind},
		{Uint16(0), Uint16Kind},
		{Uint32(0), Uint32Kind},
		{Uint64(0), Uint64Kind},
	}

	for _, d := range data {
		assert.True(t, d.v.Type().Equals(MakePrimitiveType(d.k)))
	}
}
