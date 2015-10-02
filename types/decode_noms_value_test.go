package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/attic-labs/noms/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/attic-labs/noms/chunks"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := []interface{}{int64(1), "hi", true}
	r := newJsonArrayReader(a, cs)

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
	cs := chunks.NewMemoryStore()

	a := parseJson(`[0, true]`)
	r := newJsonArrayReader(a, cs)
	k := r.readKind()
	assert.Equal(BoolKind, k)

	r = newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(BoolKind, tr.Kind())
	b := r.readValue(tr)
	assert.EqualValues(Bool(true), b)
}

func TestReadListOfInt32(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf("[%d, %d, 0, 1, 2, 3]", ListKind, Int32Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ListKind, tr.Kind())
	assert.Equal(Int32Kind, tr.Desc.(CompoundDesc).ElemTypes[0].Kind())
	l := r.readList(tr)
	assert.EqualValues(NewList(Int32(0), Int32(1), Int32(2), Int32(3)), l)
}

func TestReadListOfValue(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf(`[%d, %d, %d, 1, %d, "hi", %d, true]`, ListKind, ValueKind, Int32Kind, StringKind, BoolKind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ListKind, tr.Kind())
	assert.Equal(ValueKind, tr.Desc.(CompoundDesc).ElemTypes[0].Kind())
	l := r.readList(tr)
	assert.EqualValues(NewList(Int32(1), NewString("hi"), Bool(true)), l)
}

func TestReadValueListOfInt8(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf(`[%d, %d, %d, [0, 1, 2]]`, ValueKind, ListKind, Int8Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ValueKind, tr.Kind())
	l := r.readValue(tr)
	assert.EqualValues(NewList(Int8(0), Int8(1), Int8(2)), l)
}

func TestReadMapOfInt64ToFloat64(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf("[%d, %d, %d, 0, 1, 2, 3]", MapKind, Int64Kind, Float64Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(MapKind, tr.Kind())
	assert.Equal(Int64Kind, tr.Desc.(CompoundDesc).ElemTypes[0].Kind())
	assert.Equal(Float64Kind, tr.Desc.(CompoundDesc).ElemTypes[1].Kind())
	m := r.readMap(tr)
	assert.EqualValues(NewMap(Int64(0), Float64(1), Int64(2), Float64(3)), m)
}

func TestReadValueMapOfUInt64ToUInt32(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf("[%d, %d, %d, %d, [0, 1, 2, 3]]", ValueKind, MapKind, UInt64Kind, UInt32Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ValueKind, tr.Kind())
	m := r.readValue(tr)
	assert.True(NewMap(UInt64(0), UInt32(1), UInt64(2), UInt32(3)).Equals(m))
}

func TestReadSetOfUInt8(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf("[%d, %d, 0, 1, 2, 3]", SetKind, UInt8Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(SetKind, tr.Kind())
	assert.Equal(UInt8Kind, tr.Desc.(CompoundDesc).ElemTypes[0].Kind())
	s := r.readSet(tr)
	assert.EqualValues(NewSet(UInt8(0), UInt8(1), UInt8(2), UInt8(3)), s)
}

func TestReadValueSetOfUInt16(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	a := parseJson(fmt.Sprintf("[%d, %d, %d, [0, 1, 2, 3]]", ValueKind, SetKind, UInt16Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ValueKind, tr.Kind())
	m := r.readValue(tr)
	assert.True(NewSet(UInt16(0), UInt16(1), UInt16(2), UInt16(3)).Equals(m))
}

func TestReadStruct(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	// Cannot use parse since it is in a different package that depends on types!
	// struct A1 {
	//   x: Float32
	//   b: Bool
	//   s: String
	// }

	tref := MakeStructTypeRef("A1", []Field{
		Field{"x", MakePrimitiveTypeRef(Int16Kind), false},
		Field{"s", MakePrimitiveTypeRef(StringKind), false},
		Field{"b", MakePrimitiveTypeRef(BoolKind), false},
	}, Choices{})
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("A1", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A1", 42, "hi", true]`, StructKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(StructKind, tr.Kind())
	v := r.readStruct(tr)

	assert.True(v.Get(NewString("$name")).Equals(NewString("A1")))
	assert.True(v.Get(NewString("$type")).Equals(tref))
	assert.True(v.Get(NewString("x")).Equals(Int16(42)))
	assert.True(v.Get(NewString("s")).Equals(NewString("hi")))
	assert.True(v.Get(NewString("b")).Equals(Bool(true)))
}

func TestReadStructUnion(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	// Cannot use parse since it is in a different package that depends on types!
	// struct A2 {
	//   x: Float32
	//   union {
	//     b: Bool
	//     s: String
	//   }
	// }

	tref := MakeStructTypeRef("A2", []Field{
		Field{"x", MakePrimitiveTypeRef(Float32Kind), false},
	}, Choices{
		Field{"b", MakePrimitiveTypeRef(BoolKind), false},
		Field{"s", MakePrimitiveTypeRef(StringKind), false},
	})
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("A2", tref))
	ref := RegisterPackage(&pkg)

	a := parseJson(fmt.Sprintf(`[%d, "%s", "A2", 42, 1, "hi"]`, StructKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(StructKind, tr.Kind())
	v := r.readStruct(tr)

	assert.True(v.Get(NewString("$name")).Equals(NewString("A2")))
	assert.True(v.Get(NewString("$type")).Equals(tref))
	assert.True(v.Get(NewString("x")).Equals(Float32(42)))
	assert.False(v.Has(NewString("b")))
	assert.False(v.Has(NewString("s")))
	assert.True(v.Get(NewString("$unionIndex")).Equals(UInt32(1)))
	assert.True(v.Get(NewString("$unionValue")).Equals(NewString("hi")))
}

func TestReadStructOptional(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	// Cannot use parse since it is in a different package that depends on types!
	// struct A3 {
	//   x: Float32
	//   s: optional String
	//   b: optional Bool
	// }

	tref := MakeStructTypeRef("A3", []Field{
		Field{"x", MakePrimitiveTypeRef(Float32Kind), false},
		Field{"s", MakePrimitiveTypeRef(StringKind), true},
		Field{"b", MakePrimitiveTypeRef(BoolKind), true},
	}, Choices{})
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("A3", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A3", 42, false, true, false]`, StructKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(StructKind, tr.Kind())
	v := r.readStruct(tr)

	assert.True(v.Get(NewString("$name")).Equals(NewString("A3")))
	assert.True(v.Get(NewString("$type")).Equals(tref))
	assert.True(v.Get(NewString("x")).Equals(Float32(42)))
	assert.False(v.Has(NewString("s")))
	assert.True(v.Get(NewString("b")).Equals(Bool(false)))
}

func TestReadStructWithList(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	// Cannot use parse since it is in a different package that depends on types!
	// struct A4 {
	//   b: Bool
	//   l: List(Int32)
	//   s: String
	// }

	tref := MakeStructTypeRef("A4", []Field{
		Field{"b", MakePrimitiveTypeRef(BoolKind), false},
		Field{"l", MakeCompoundTypeRef("", ListKind, MakePrimitiveTypeRef(Int32Kind)), false},
		Field{"s", MakePrimitiveTypeRef(StringKind), false},
	}, Choices{})
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("A4", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A4", true, [0, 1, 2], "hi"]`, StructKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(StructKind, tr.Kind())
	v := r.readStruct(tr)

	assert.True(v.Get(NewString("$name")).Equals(NewString("A4")))
	assert.True(v.Get(NewString("$type")).Equals(tref))
	assert.True(v.Get(NewString("b")).Equals(Bool(true)))
	assert.True(v.Get(NewString("l")).Equals(NewList(Int32(0), Int32(1), Int32(2))))
	assert.True(v.Get(NewString("s")).Equals(NewString("hi")))
}

func TestReadStructWithValue(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	// Cannot use parse since it is in a different package that depends on types!
	// struct A5 {
	//   b: Bool
	//   v: Value
	//   s: String
	// }

	tref := MakeStructTypeRef("A5", []Field{
		Field{"b", MakePrimitiveTypeRef(BoolKind), false},
		Field{"v", MakePrimitiveTypeRef(ValueKind), false},
		Field{"s", MakePrimitiveTypeRef(StringKind), false},
	}, Choices{})
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("A5", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "A5", true, %d, 42, "hi"]`, StructKind, ref.String(), UInt8Kind))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(StructKind, tr.Kind())
	v := r.readStruct(tr)

	assert.True(v.Get(NewString("$name")).Equals(NewString("A5")))
	assert.True(v.Get(NewString("$type")).Equals(tref))
	assert.True(v.Get(NewString("b")).Equals(Bool(true)))
	assert.True(v.Get(NewString("v")).Equals(UInt8(42)))
	assert.True(v.Get(NewString("s")).Equals(NewString("hi")))
}

func TestReadValueStruct(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	// Cannot use parse since it is in a different package that depends on types!
	// struct A1 {
	//   x: Float32
	//   b: Bool
	//   s: String
	// }

	tref := MakeStructTypeRef("A1", []Field{
		Field{"x", MakePrimitiveTypeRef(Int16Kind), false},
		Field{"s", MakePrimitiveTypeRef(StringKind), false},
		Field{"b", MakePrimitiveTypeRef(BoolKind), false},
	}, Choices{})
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("A1", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, %d, "%s", "A1", 42, "hi", true]`, ValueKind, StructKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ValueKind, tr.Kind())
	v := r.readValue(tr).(Map)

	assert.True(v.Get(NewString("$name")).Equals(NewString("A1")))
	assert.True(v.Get(NewString("$type")).Equals(tref))
	assert.True(v.Get(NewString("x")).Equals(Int16(42)))
	assert.True(v.Get(NewString("s")).Equals(NewString("hi")))
	assert.True(v.Get(NewString("b")).Equals(Bool(true)))
}

func TestReadEnum(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	tref := MakeEnumTypeRef("E", "a", "b", "c")
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("E", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, "%s", "E", 1]`, EnumKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(EnumKind, tr.Kind())
	v := r.readEnum(tr)
	assert.Equal(uint32(1), uint32(v.(UInt32)))
}

func TestReadValueEnum(t *testing.T) {
	assert := assert.New(t)
	cs := chunks.NewMemoryStore()

	tref := MakeEnumTypeRef("E", "a", "b", "c")
	pkg := NewPackage().SetNamedTypes(NewMapOfStringToTypeRef().Set("E", tref))
	ref := RegisterPackage(&pkg)

	// TODO: Should use ordinal of type and not name
	a := parseJson(fmt.Sprintf(`[%d, %d, "%s", "E", 1]`, ValueKind, EnumKind, ref.String()))
	r := newJsonArrayReader(a, cs)
	tr := r.readTypeRef()
	assert.Equal(ValueKind, tr.Kind())
	v := r.readValue(tr)
	assert.Equal(uint32(1), uint32(v.(UInt32)))
}
