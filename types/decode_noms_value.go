package types

import (
	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/ref"
)

type valueAsNomsValue struct {
	v Value
}

func (v valueAsNomsValue) NomsValue() Value {
	return v.v
}

func (v valueAsNomsValue) TypeRef() TypeRef {
	return v.TypeRef()
}

func fromTypedEncodeable(w typedValueWrapper, cs chunks.ChunkSource) NomsValue {
	i := w.TypedValue()
	r := newJsonArrayReader(i, cs)
	t := r.readTypeRef()
	v := r.readTopLevelValue(t)
	switch v := v.(type) {
	case NomsValue:
		return v
	case Value:
		return valueAsNomsValue{v}
	}
	panic("unreachable")
}

type jsonArrayReader struct {
	a  []interface{}
	i  int
	cs chunks.ChunkSource
}

func newJsonArrayReader(a []interface{}, cs chunks.ChunkSource) *jsonArrayReader {
	return &jsonArrayReader{a: a, i: 0, cs: cs}
}

func (r *jsonArrayReader) read() interface{} {
	v := r.a[r.i]
	r.i++
	return v
}

func (r *jsonArrayReader) atEnd() bool {
	return r.i >= len(r.a)
}

func (r *jsonArrayReader) readString() string {
	return r.read().(string)
}

func (r *jsonArrayReader) readInt64() int64 {
	return r.read().(int64)
}

func (r *jsonArrayReader) readBool() bool {
	return r.read().(bool)
}

func (r *jsonArrayReader) readArray() []interface{} {
	return r.read().([]interface{})
}

func (r *jsonArrayReader) readKind() NomsKind {
	return NomsKind(r.read().(float64))
}

func (r *jsonArrayReader) readRef() ref.Ref {
	s := r.readString()
	return ref.Parse(s)
}

func (r *jsonArrayReader) readPackage() *Package {
	ref := r.readRef()
	// TODO: Should load the package if not registered?
	return LookupPackage(ref)
}

func (r *jsonArrayReader) readTypeRef() TypeRef {
	kind := r.readKind()
	if IsPrimitiveKind(kind) {
		return MakePrimitiveTypeRef(kind)
	}
	switch kind {
	case ListKind, SetKind, RefKind:
		elemType := r.readTypeRef()
		return MakeCompoundTypeRef("", kind, elemType)
	case MapKind:
		keyType := r.readTypeRef()
		valueType := r.readTypeRef()
		return MakeCompoundTypeRef("", kind, keyType, valueType)
	case EnumKind, StructKind:
		pkg := r.readPackage()
		d.Chk.NotNil(pkg)
		name := r.readString()
		// TODO: This is wrong. Package uses name to TypeRef
		return pkg.NamedTypes().Get(name)
	}
	panic("unreachable")
}

func (r *jsonArrayReader) readList(t TypeRef) List {
	desc := t.Desc.(CompoundDesc)
	ll := []Value{}
	elemType := desc.ElemTypes[0]
	for !r.atEnd() {
		v := r.readValue(elemType)
		ll = append(ll, v)
	}
	return NewList(ll...)
}

func (r *jsonArrayReader) readSet(t TypeRef) Set {
	desc := t.Desc.(CompoundDesc)
	ll := []Value{}
	elemType := desc.ElemTypes[0]
	for !r.atEnd() {
		v := r.readValue(elemType)
		ll = append(ll, v)
	}
	return NewSet(ll...)
}

func (r *jsonArrayReader) readMap(t TypeRef) Map {
	desc := t.Desc.(CompoundDesc)
	ll := []Value{}
	keyType := desc.ElemTypes[0]
	valueType := desc.ElemTypes[1]
	for !r.atEnd() {
		k := r.readValue(keyType)
		v := r.readValue(valueType)
		ll = append(ll, k, v)
	}
	return NewMap(ll...)
}

func (r *jsonArrayReader) readStruct(t TypeRef) Map {
	desc := t.Desc.(StructDesc)
	m := NewMap(
		NewString("$name"), NewString(t.Name()),
		NewString("$type"), t)

	for _, f := range desc.Fields {
		if f.Optional {
			b := r.read().(bool)
			if b {
				v := r.readValue(f.T)
				m = m.Set(NewString(f.Name), v)
			}
		} else {
			v := r.readValue(f.T)
			m = m.Set(NewString(f.Name), v)
		}
	}
	if len(desc.Union) > 0 {
		i := uint32(r.read().(float64))
		m = m.Set(NewString("$unionIndex"), UInt32(i))
		v := r.readValue(desc.Union[i].T)
		m = m.Set(NewString("$unionValue"), v)
	}

	return m
}

func (r *jsonArrayReader) readEnum(TypeRef) Value {
	return UInt32(r.read().(float64))
}

func (r *jsonArrayReader) readValue(t TypeRef) Value {
	switch t.Kind() {
	case ListKind, MapKind, SetKind:
		a := r.readArray()
		r2 := newJsonArrayReader(a, r.cs)
		return r2.readTopLevelValue(t)
	default:
		return r.readTopLevelValue(t)
	}
}

func (r *jsonArrayReader) readTopLevelValue(t TypeRef) Value {
	switch t.Kind() {
	case BoolKind:
		return Bool(r.read().(bool))
	case UInt8Kind:
		return UInt8(r.read().(float64))
	case UInt16Kind:
		return UInt16(r.read().(float64))
	case UInt32Kind:
		return UInt32(r.read().(float64))
	case UInt64Kind:
		return UInt64(r.read().(float64))
	case Int8Kind:
		return Int8(r.read().(float64))
	case Int16Kind:
		return Int16(r.read().(float64))
	case Int32Kind:
		return Int32(r.read().(float64))
	case Int64Kind:
		return Int64(r.read().(float64))
	case Float32Kind:
		return Float32(r.read().(float64))
	case Float64Kind:
		return Float64(r.read().(float64))
	case StringKind:
		return NewString(r.readString())
	case BlobKind:
		panic("not implemented")
	case ValueKind:
		// The value is always tagged
		t := r.readTypeRef()
		return r.readValue(t)
	case ListKind:
		return r.readList(t)
	case MapKind:
		return r.readMap(t)
	case RefKind:
		panic("not implemented")
	case SetKind:
		return r.readSet(t)
	case EnumKind:
		return r.readEnum(t)
	case StructKind:
		return r.readStruct(t)
	case TypeRefKind:
		panic("not implemented")
	}
	panic("not reachable")
}
