package csv

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/datas"
	"github.com/attic-labs/noms/types"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)
	ds := datas.NewDataStore(chunks.NewMemoryStore())

	dataString := `a,1,true
b,2,false
`
	r := NewCSVReader(bytes.NewBufferString(dataString), ',')

	headers := []string{"A", "B", "C"}
	kinds := KindSlice{types.StringKind, types.Int8Kind, types.BoolKind}
	l, typeRef, typeDef := Read(r, "test", headers, kinds, ds)

	assert.Equal(uint64(2), l.Len())

	assert.True(typeRef.IsUnresolved())

	desc, ok := typeDef.Desc.(types.StructDesc)
	assert.True(ok)
	assert.Len(desc.Fields, 3)
	assert.Equal("A", desc.Fields[0].Name)
	assert.Equal("B", desc.Fields[1].Name)
	assert.Equal("C", desc.Fields[2].Name)

	assert.True(l.Get(0).(types.Struct).Get("A").Equals(types.NewString("a")))
	assert.True(l.Get(1).(types.Struct).Get("A").Equals(types.NewString("b")))

	assert.True(l.Get(0).(types.Struct).Get("B").Equals(types.Int8(1)))
	assert.True(l.Get(1).(types.Struct).Get("B").Equals(types.Int8(2)))

	assert.True(l.Get(0).(types.Struct).Get("C").Equals(types.Bool(true)))
	assert.True(l.Get(1).(types.Struct).Get("C").Equals(types.Bool(false)))
}

func testTrailingHelper(t *testing.T, dataString string) {
	assert := assert.New(t)
	ds := datas.NewDataStore(chunks.NewMemoryStore())

	r := NewCSVReader(bytes.NewBufferString(dataString), ',')

	headers := []string{"A", "B"}
	kinds := KindSlice{types.StringKind, types.StringKind}
	l, typeRef, typeDef := Read(r, "test", headers, kinds, ds)

	assert.Equal(uint64(3), l.Len())

	assert.True(typeRef.IsUnresolved())

	desc, ok := typeDef.Desc.(types.StructDesc)
	assert.True(ok)
	assert.Len(desc.Fields, 2)
	assert.Equal("A", desc.Fields[0].Name)
	assert.Equal("B", desc.Fields[1].Name)
}

func TestReadTrailingHole(t *testing.T) {
	dataString := `a,b,
d,e,
g,h,
`
	testTrailingHelper(t, dataString)
}

func TestReadTrailingHoles(t *testing.T) {
	dataString := `a,b,,
d,e
g,h
`
	testTrailingHelper(t, dataString)
}

func TestReadTrailingValues(t *testing.T) {
	dataString := `a,b
d,e,f
g,h,i,j
`
	testTrailingHelper(t, dataString)
}

func TestReadParseError(t *testing.T) {
	assert := assert.New(t)
	ds := datas.NewDataStore(chunks.NewMemoryStore())

	dataString := `a,"b`
	r := NewCSVReader(bytes.NewBufferString(dataString), ',')

	headers := []string{"A", "B"}
	kinds := KindSlice{types.StringKind, types.StringKind}
	func() {
		defer func() {
			r := recover()
			assert.NotNil(r)
			_, ok := r.(*csv.ParseError)
			assert.True(ok, "Should be a ParseError")
		}()
		Read(r, "test", headers, kinds, ds)
	}()
}
