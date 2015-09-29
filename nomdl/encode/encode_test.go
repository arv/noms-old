package encode

import (
	"testing"

	"github.com/attic-labs/noms/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

func TestWrite(t *testing.T) {
	assert := assert.New(t)

	w := newJsonArrayWriter()
	w.writeTypeRef(types.MakePrimitiveTypeRef(types.UInt64Kind))

	assert.EqualValues([]interface{}{types.UInt64Kind}, *w)
}

func TestWriteValue(t *testing.T) {
	assert := assert.New(t)

	f := func(k types.NomsKind, v types.Value, ex interface{}) {
		w := newJsonArrayWriter()
		tref := types.MakePrimitiveTypeRef(k)
		w.writeTypeRef(tref)
		w.writeValue(tref, v)
		assert.EqualValues([]interface{}{k, ex}, *w)
	}

	f(types.BoolKind, types.Bool(true), true)
	f(types.BoolKind, types.Bool(false), false)

	f(types.UInt8Kind, types.UInt8(0), uint8(0))
	f(types.UInt16Kind, types.UInt16(0), uint16(0))
	f(types.UInt32Kind, types.UInt32(0), uint32(0))
	f(types.UInt64Kind, types.UInt64(0), uint64(0))
	f(types.Int8Kind, types.Int8(0), int8(0))
	f(types.Int16Kind, types.Int16(0), int16(0))
	f(types.Int32Kind, types.Int32(0), int32(0))
	f(types.Int64Kind, types.Int64(0), int64(0))
	f(types.Float32Kind, types.Float32(0), float32(0))
	f(types.Float64Kind, types.Float64(0), float64(0))

	f(types.StringKind, types.NewString("hi"), "hi")
}

func TestWriteList(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.Int32Kind))
	v := types.NewList(types.Int32(0), types.Int32(1), types.Int32(2), types.Int32(3))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeList(tref, v)
	assert.EqualValues([]interface{}{types.ListKind, types.Int32Kind, int32(0), int32(1), int32(2), int32(3)}, *w)
}

func TestWriteListOfList(t *testing.T) {
	assert := assert.New(t)

	it := types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.Int16Kind))
	tref := types.MakeCompoundTypeRef("", types.ListKind, it)
	v := types.NewList(types.NewList(types.Int16(0)), types.NewList(types.Int16(1), types.Int16(2), types.Int16(3)))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeList(tref, v)
	assert.EqualValues([]interface{}{types.ListKind, types.ListKind, types.Int16Kind, []interface{}{int16(0)}, []interface{}{int16(1), int16(2), int16(3)}}, *w)
}

func TestWriteSet(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeCompoundTypeRef("", types.SetKind, types.MakePrimitiveTypeRef(types.UInt32Kind))
	v := types.NewSet(types.UInt32(0), types.UInt32(1), types.UInt32(2), types.UInt32(3))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeSet(tref, v)
	// the order of the elements is based on the ref of the value.
	assert.EqualValues([]interface{}{types.SetKind, types.UInt32Kind, uint32(3), uint32(1), uint32(0), uint32(2)}, *w)
}

func TestWriteSetOfSet(t *testing.T) {
	assert := assert.New(t)

	st := types.MakeCompoundTypeRef("", types.SetKind, types.MakePrimitiveTypeRef(types.Int32Kind))
	tref := types.MakeCompoundTypeRef("", types.SetKind, st)
	v := types.NewSet(types.NewSet(types.Int32(0)), types.NewSet(types.Int32(1), types.Int32(2), types.Int32(3)))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeSet(tref, v)
	// the order of the elements is based on the ref of the value.
	assert.EqualValues([]interface{}{types.SetKind, types.SetKind, types.Int32Kind, []interface{}{int32(0)}, []interface{}{int32(1), int32(3), int32(2)}}, *w)
}

func TestWriteMap(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeCompoundTypeRef("", types.MapKind, types.MakePrimitiveTypeRef(types.StringKind), types.MakePrimitiveTypeRef(types.BoolKind))
	v := types.NewMap(types.NewString("a"), types.Bool(false), types.NewString("b"), types.Bool(true))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeMap(tref, v)
	// the order of the elements is based on the ref of the value.
	assert.EqualValues([]interface{}{types.MapKind, types.StringKind, types.BoolKind, "a", false, "b", true}, *w)
}

func TestWriteMapOfMap(t *testing.T) {
	assert := assert.New(t)

	kt := types.MakeCompoundTypeRef("", types.MapKind, types.MakePrimitiveTypeRef(types.StringKind), types.MakePrimitiveTypeRef(types.Int64Kind))
	vt := types.MakeCompoundTypeRef("", types.SetKind, types.MakePrimitiveTypeRef(types.BoolKind))
	tref := types.MakeCompoundTypeRef("", types.MapKind, kt, vt)
	v := types.NewMap(types.NewMap(types.NewString("a"), types.Int64(0)), types.NewSet(types.Bool(true)))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeMap(tref, v)
	// the order of the elements is based on the ref of the value.
	assert.EqualValues([]interface{}{types.MapKind, types.MapKind, types.StringKind, types.Int64Kind, types.SetKind, types.BoolKind, []interface{}{"a", int64(0)}, []interface{}{true}}, *w)
}

func TestWriteEmptyStruct(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeStructTypeRef("S", []types.Field{}, types.Choices{})
	v := types.NewMap()

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S"}, *w)
}

func TestWriteStruct(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeStructTypeRef("S", []types.Field{
		types.Field{"x", types.MakePrimitiveTypeRef(types.Int8Kind), false},
		types.Field{"b", types.MakePrimitiveTypeRef(types.BoolKind), false},
	}, types.Choices{})
	v := types.NewMap(types.NewString("x"), types.Int8(42), types.NewString("b"), types.Bool(true))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", int8(42), true}, *w)
}

func TestWriteStructOptionalField(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeStructTypeRef("S", []types.Field{
		types.Field{"x", types.MakePrimitiveTypeRef(types.Int8Kind), true},
		types.Field{"b", types.MakePrimitiveTypeRef(types.BoolKind), false},
	}, types.Choices{})
	v := types.NewMap(types.NewString("x"), types.Int8(42), types.NewString("b"), types.Bool(true))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", uint32(1), int8(42), true}, *w)

	v = types.NewMap(types.NewString("b"), types.Bool(true))

	w = newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", uint32(0), uint32(0), true}, *w)
}

func TestWriteStructWithUnion(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeStructTypeRef("S", []types.Field{
		types.Field{"x", types.MakePrimitiveTypeRef(types.Int8Kind), false},
	}, types.Choices{
		types.Field{"b", types.MakePrimitiveTypeRef(types.BoolKind), false},
		types.Field{"s", types.MakePrimitiveTypeRef(types.StringKind), false},
	})
	v := types.NewMap(types.NewString("x"), types.Int8(42), types.NewString("$unionIndex"), types.UInt32(1), types.NewString("$unionValue"), types.NewString("hi"))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", int8(42), uint32(1), "hi"}, *w)

	v = types.NewMap(types.NewString("x"), types.Int8(42), types.NewString("$unionIndex"), types.UInt32(0), types.NewString("$unionValue"), types.Bool(true))

	w = newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", int8(42), uint32(0), true}, *w)
}

func TestWriteStructWithList(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeStructTypeRef("S", []types.Field{
		types.Field{"l", types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.StringKind)), false},
	}, types.Choices{})
	v := types.NewMap(types.NewString("l"), types.NewList(types.NewString("a"), types.NewString("b")))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", []interface{}{"a", "b"}}, *w)

	v = types.NewMap(types.NewString("l"), types.NewList())
	w = newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", []interface{}{}}, *w)
}

func TestWriteStructWithStruc(t *testing.T) {
	assert := assert.New(t)

	st := types.MakeStructTypeRef("S2", []types.Field{
		types.Field{"x", types.MakePrimitiveTypeRef(types.Int32Kind), false},
	}, types.Choices{})
	tref := types.MakeStructTypeRef("S", []types.Field{
		types.Field{"s", st, false},
	}, types.Choices{})
	v := types.NewMap(types.NewString("s"), types.NewMap(types.NewString("x"), types.Int32(42)))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeStruct(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.StructKind, ref.String(), "S", int32(42)}, *w)
}

func TestWriteEnum(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeEnumTypeRef("E", "a", "b", "c")
	v := types.UInt32(1)

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeEnum(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.EnumKind, ref.String(), "E", uint32(1)}, *w)
}

func TestWriteListOfEnum(t *testing.T) {
	assert := assert.New(t)

	et := types.MakeEnumTypeRef("E", "a", "b", "c")
	tref := types.MakeCompoundTypeRef("", types.ListKind, et)
	v := types.NewList(types.UInt32(0), types.UInt32(1), types.UInt32(2))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeList(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.ListKind, types.EnumKind, ref.String(), "E", uint32(0), uint32(1), uint32(2)}, *w)
}

func TestWriteListOfValue(t *testing.T) {
	assert := assert.New(t)

	tref := types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.ValueKind))
	v := types.NewList(
		types.Bool(true),
		types.UInt8(1),
		types.UInt16(1),
		types.UInt32(1),
		types.UInt64(1),
		types.Int8(1),
		types.Int16(1),
		types.Int32(1),
		types.Int64(1),
		types.Float32(1),
		types.Float64(1),
		types.NewString("hi"),
	)

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeList(tref, v)
	assert.EqualValues([]interface{}{types.ListKind, types.ValueKind,
		types.BoolKind, true,
		types.UInt8Kind, uint8(1),
		types.UInt16Kind, uint16(1),
		types.UInt32Kind, uint32(1),
		types.UInt64Kind, uint64(1),
		types.Int8Kind, int8(1),
		types.Int16Kind, int16(1),
		types.Int32Kind, int32(1),
		types.Int64Kind, int64(1),
		types.Float32Kind, float32(1),
		types.Float64Kind, float64(1),
		types.StringKind, "hi",
	}, *w)
}

func TestWriteListOfValueWithStruct(t *testing.T) {
	assert := assert.New(t)
	tref := types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.ValueKind))
	st := types.MakeStructTypeRef("S", []types.Field{
		types.Field{"x", types.MakePrimitiveTypeRef(types.Int32Kind), false},
	}, types.Choices{})
	v := types.NewList(types.NewMap(types.NewString("$type"), st, types.NewString("x"), types.Int32(42)))

	w := newJsonArrayWriter()
	w.writeTypeRef(tref)
	w.writeList(tref, v)
	ref := ref.Ref{}
	assert.EqualValues([]interface{}{types.ListKind, types.ValueKind, types.StructKind, ref.String(), "S", int32(42)}, *w)
}
