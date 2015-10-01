package decode

import (
	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

type jsonArrayReader struct {
	a []interface{}
	i int
}

func newJsonArrayReader(a []interface{}) *jsonArrayReader {
	return &jsonArrayReader{a: a, i: 0}
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

func (r *jsonArrayReader) readKind() types.NomsKind {
	return types.NomsKind(r.read().(float64))
}

func (r *jsonArrayReader) readRef() ref.Ref {
	s := r.readString()
	return ref.Parse(s)
}

func (r *jsonArrayReader) readPackage(cs chunks.ChunkSource) *types.Package {
	ref := r.readRef()
	// TODO: Should load the package if not registered?
	return types.LookupPackage(ref)
}

func (r *jsonArrayReader) readTypeRef() types.TypeRef {
	kind := r.readKind()
	if types.IsPrimitiveKind(kind) {
		return types.MakePrimitiveTypeRef(kind)
	}
	switch kind {
	case types.ListKind, types.SetKind, types.RefKind:
		elemType := r.readTypeRef()
		return types.MakeCompoundTypeRef("", kind, elemType)
	case types.MapKind:
		keyType := r.readTypeRef()
		valueType := r.readTypeRef()
		return types.MakeCompoundTypeRef("", kind, keyType, valueType)
	case types.EnumKind, types.StructKind:
		// TODO: Provide chunk source?
		pkg := r.readPackage(nil)
		name := r.readString()
		// TODO: This is wrong. Package uses name to TypeRef
		return pkg.NamedTypes().Get(name)
	}
	panic("unreachable")
}

func (r *jsonArrayReader) readList(t types.TypeRef) types.List {
	desc := t.Desc.(types.CompoundDesc)
	ll := []types.Value{}
	elemType := desc.ElemTypes[0]
	for !r.atEnd() {
		v := r.readValue(elemType)
		ll = append(ll, v)
	}
	return types.NewList(ll...)
}

func (r *jsonArrayReader) readSet(t types.TypeRef) types.Set {
	desc := t.Desc.(types.CompoundDesc)
	ll := []types.Value{}
	elemType := desc.ElemTypes[0]
	for !r.atEnd() {
		v := r.readValue(elemType)
		ll = append(ll, v)
	}
	return types.NewSet(ll...)
}

func (r *jsonArrayReader) readMap(t types.TypeRef) types.Map {
	desc := t.Desc.(types.CompoundDesc)
	ll := []types.Value{}
	keyType := desc.ElemTypes[0]
	valueType := desc.ElemTypes[1]
	for !r.atEnd() {
		k := r.readValue(keyType)
		v := r.readValue(valueType)
		ll = append(ll, k, v)
	}
	return types.NewMap(ll...)
}

func (r *jsonArrayReader) readStruct(t types.TypeRef) types.Map {
	desc := t.Desc.(types.StructDesc)
	m := types.NewMap(
		types.NewString("$name"), types.NewString(t.Name()),
		types.NewString("$type"), t)

	for _, f := range desc.Fields {
		if f.Optional {
			p := r.read().(float64)
			if p != 0 {
				v := r.readValue(f.T)
				m = m.Set(types.NewString(f.Name), v)
			}
		} else {
			v := r.readValue(f.T)
			m = m.Set(types.NewString(f.Name), v)
		}
	}
	if len(desc.Union) > 0 {
		i := uint32(r.read().(float64))
		m = m.Set(types.NewString("$unionIndex"), types.UInt32(i))
		v := r.readValue(desc.Union[i].T)
		m = m.Set(types.NewString("$unionValue"), v)
	}

	return m
}

func (r *jsonArrayReader) readEnum(types.TypeRef) types.Value {
	return types.UInt32(r.read().(float64))
}

func (r *jsonArrayReader) readValue(t types.TypeRef) types.Value {
	switch t.Kind() {
	case types.BoolKind:
		return types.Bool(r.read().(bool))
	case types.UInt8Kind:
		return types.UInt8(r.read().(float64))
	case types.UInt16Kind:
		return types.UInt16(r.read().(float64))
	case types.UInt32Kind:
		return types.UInt32(r.read().(float64))
	case types.UInt64Kind:
		return types.UInt64(r.read().(float64))
	case types.Int8Kind:
		return types.Int8(r.read().(float64))
	case types.Int16Kind:
		return types.Int16(r.read().(float64))
	case types.Int32Kind:
		return types.Int32(r.read().(float64))
	case types.Int64Kind:
		return types.Int64(r.read().(float64))
	case types.Float32Kind:
		return types.Float32(r.read().(float64))
	case types.Float64Kind:
		return types.Float64(r.read().(float64))
	case types.StringKind:
		return types.NewString(r.readString())
	case types.BlobKind:
		panic("not implemented")
	case types.ValueKind:
		// The value is always tagged
		t := r.readTypeRef()
		return r.readValue(t)
	case types.ListKind:
		a := r.readArray()
		r2 := newJsonArrayReader(a)
		return r2.readList(t)
	case types.MapKind:
		a := r.readArray()
		r2 := newJsonArrayReader(a)
		return r2.readMap(t)
	case types.RefKind:
		panic("not implemented")
	case types.SetKind:
		a := r.readArray()
		r2 := newJsonArrayReader(a)
		return r2.readSet(t)
	case types.EnumKind:
		return r.readEnum(t)
	case types.StructKind:
		return r.readStruct(t)
	case types.TypeRefKind:
		panic("not implemented")
	}
	panic("not reachable")
}
