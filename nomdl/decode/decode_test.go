package decode

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/attic-labs/noms/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/attic-labs/noms/types"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)

	a := []interface{}{int64(1), "hi", true}
	r := newJsonArrayReader(a)

	assert.Equal(int64(1), r.readInt64())
	assert.False(r.atEnd())

	assert.Equal("hi", r.readString())
	assert.False(r.atEnd())

	assert.Equal(true, r.readBool())
	assert.True(r.atEnd())
}

func parseJson(s string) (v []interface{}) {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.Decode(&v)
	return
}

func TestReadTypeRef(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(`[0, true]`)
	r := newJsonArrayReader(a)
	k := r.readKind()
	assert.Equal(types.BoolKind, k)

	r = newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.BoolKind, tr.Kind())
	b := r.readValue(tr)
	assert.EqualValues(types.Bool(true), b)
}

func TestReadListOfInt32(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf("[%d, %d, 0, 1, 2, 3]", types.ListKind, types.Int32Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ListKind, tr.Kind())
	assert.Equal(types.Int32Kind, tr.Desc.(types.CompoundDesc).ElemTypes[0].Kind())
	l := r.readList(tr)
	assert.EqualValues(types.NewList(types.Int32(0), types.Int32(1), types.Int32(2), types.Int32(3)), l)
}

func TestReadListOfValue(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf(`[%d, %d, %d, 1, %d, "hi", %d, true]`, types.ListKind, types.ValueKind, types.Int32Kind, types.StringKind, types.BoolKind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ListKind, tr.Kind())
	assert.Equal(types.ValueKind, tr.Desc.(types.CompoundDesc).ElemTypes[0].Kind())
	l := r.readList(tr)
	assert.EqualValues(types.NewList(types.Int32(1), types.NewString("hi"), types.Bool(true)), l)
}

func TestReadValueListOfInt8(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf(`[%d, %d, %d, [0, 1, 2]]`, types.ValueKind, types.ListKind, types.Int8Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ValueKind, tr.Kind())
	l := r.readValue(tr)
	assert.EqualValues(types.NewList(types.Int8(0), types.Int8(1), types.Int8(2)), l)
}

func TestReadMapOfInt64ToFloat64(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf("[%d, %d, %d, 0, 1, 2, 3]", types.MapKind, types.Int64Kind, types.Float64Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.MapKind, tr.Kind())
	assert.Equal(types.Int64Kind, tr.Desc.(types.CompoundDesc).ElemTypes[0].Kind())
	assert.Equal(types.Float64Kind, tr.Desc.(types.CompoundDesc).ElemTypes[1].Kind())
	m := r.readMap(tr)
	assert.EqualValues(types.NewMap(types.Int64(0), types.Float64(1), types.Int64(2), types.Float64(3)), m)
}

func TestReadValueMapOfUInt64ToUInt32(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf("[%d, %d, %d, %d, [0, 1, 2, 3]]", types.ValueKind, types.MapKind, types.UInt64Kind, types.UInt32Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ValueKind, tr.Kind())
	m := r.readValue(tr)
	assert.True(types.NewMap(types.UInt64(0), types.UInt32(1), types.UInt64(2), types.UInt32(3)).Equals(m))
}

func TestReadSetOfUInt8(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf("[%d, %d, 0, 1, 2, 3]", types.SetKind, types.UInt8Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.SetKind, tr.Kind())
	assert.Equal(types.UInt8Kind, tr.Desc.(types.CompoundDesc).ElemTypes[0].Kind())
	s := r.readSet(tr)
	assert.EqualValues(types.NewSet(types.UInt8(0), types.UInt8(1), types.UInt8(2), types.UInt8(3)), s)
}

func TestReadValueSetOfUInt16(t *testing.T) {
	assert := assert.New(t)

	a := parseJson(fmt.Sprintf("[%d, %d, %d, [0, 1, 2, 3]]", types.ValueKind, types.SetKind, types.UInt16Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ValueKind, tr.Kind())
	m := r.readValue(tr)
	assert.True(types.NewSet(types.UInt16(0), types.UInt16(1), types.UInt16(2), types.UInt16(3)).Equals(m))
}

func TestReadStruct(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A1", 42, "hi", true]`, types.StructKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.StructKind, tr.Kind())
	v := r.readStruct(tr)

	// TODO: Package ref are currently not propagated in the parser.
	s := A1FromVal(v)
	assert.Equal(s.Def(), A1Def{42, "hi", true})
}

func TestReadStructUnion(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A2", 42, 1, "hi"]`, types.StructKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.StructKind, tr.Kind())
	v := r.readStruct(tr)

	s := A2FromVal(v)
	assert.Equal(float32(42), s.X())
	_, ok := s.B()
	assert.False(ok)
	str, ok := s.S()
	assert.True(ok)
	assert.Equal("hi", str)
}

func TestReadStructOptional(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A3", 42, 0, 0, 1, false]`, types.StructKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.StructKind, tr.Kind())
	v := r.readStruct(tr)

	s := A3FromVal(v)
	assert.Equal(float32(42), s.X())
	_, ok := s.S()
	assert.False(ok)
	b, ok := s.B()
	assert.True(ok)
	assert.Equal(false, b)
}

func TestReadStructWithList(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A4", true, [0, 1, 2], "hi"]`, types.StructKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.StructKind, tr.Kind())
	v := r.readStruct(tr)

	// TODO: Package ref are currently not propagated in the parser.
	s := A4FromVal(v)
	assert.Equal(s.Def(), A4Def{true, []int32{0, 1, 2}, "hi"})
}

func TestReadStructWithValue(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A5", true, %d, 42, "hi"]`, types.StructKind, ref.String(), types.UInt8Kind))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.StructKind, tr.Kind())
	v := r.readStruct(tr)

	// TODO: Package ref are currently not propagated in the parser.
	s := A5FromVal(v)
	assert.Equal(s.Def(), A5Def{true, types.UInt8(42), "hi"})
}

func TestReadValueStruct(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, %d, "%s", "A1", 42, "hi", true]`, types.ValueKind, types.StructKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ValueKind, tr.Kind())
	v := r.readValue(tr)
	s := A1FromVal(v)
	assert.Equal(int16(42), s.X())
	assert.Equal(true, s.B())
	assert.Equal("hi", s.S())
}

func TestReadEnum(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "E", 1]`, types.EnumKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.EnumKind, tr.Kind())
	v := r.readEnum(tr)
	assert.Equal(uint32(1), uint32(v.(types.UInt32)))
}

func TestReadValueEnum(t *testing.T) {
	assert := assert.New(t)

	ref := __decodePackageInFile_types_CachedRef
	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, %d, "%s", "E", 1]`, types.ValueKind, types.EnumKind, ref.String()))
	r := newJsonArrayReader(a)
	tr := r.readTypeRef()
	assert.Equal(types.ValueKind, tr.Kind())
	v := r.readValue(tr)
	assert.Equal(uint32(1), uint32(v.(types.UInt32)))
}
