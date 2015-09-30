package encode

import (
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

type jsonArrayWriter []interface{}

func newJsonArrayWriter() *jsonArrayWriter {
	return &jsonArrayWriter{}
}

func (w *jsonArrayWriter) write(v interface{}) {
	*w = append(*w, v)
}

func (w *jsonArrayWriter) toArray() []interface{} {
	return *w
}

func (w *jsonArrayWriter) writeRef(r ref.Ref) {
	w.write(r.String())
}

func (w *jsonArrayWriter) writeTypeRef(t types.TypeRef) {
	// TODO: Resolve if needed
	k := t.Kind()
	w.write(k)
	switch k {
	case types.EnumKind, types.StructKind:
		w.writeRef(t.PackageRef())
		// TODO: Should be ordinal instead of name.
		w.write(t.Name())
	case types.ListKind, types.MapKind, types.RefKind, types.SetKind:
		for _, elemType := range t.Desc.(types.CompoundDesc).ElemTypes {
			w.writeTypeRef(elemType)
		}
	}
}

type encodeableValue interface {
	TypeRef() types.TypeRef
}

type primitive interface {
	ToPrimitive() interface{}
}

type nomsValue interface {
	TypeRef() types.TypeRef
	NomsValue() types.Value
}

func (w *jsonArrayWriter) writeNomsValue(nv nomsValue) {
	v := nv.NomsValue()
	t := nv.TypeRef()
	w.writeTypeRef(t)
	w.writeTopLevelValue(t, v)
}

func (w *jsonArrayWriter) writeValue(t types.TypeRef, v encodeableValue) {
	switch t.Kind() {
	case types.ListKind, types.MapKind, types.SetKind:
		w2 := newJsonArrayWriter()
		w2.writeTopLevelValue(t, v)
		w.write(w2.toArray())
	default:
		w.writeTopLevelValue(t, v)
	}
}

func (w *jsonArrayWriter) writeTopLevelValue(t types.TypeRef, v encodeableValue) {
	switch t.Kind() {
	case types.BoolKind, types.Float32Kind, types.Float64Kind, types.Int16Kind, types.Int32Kind, types.Int64Kind, types.Int8Kind, types.UInt16Kind, types.UInt32Kind, types.UInt64Kind, types.UInt8Kind:
		w.write(v.(primitive).ToPrimitive())
	case types.StringKind:
		w.write(v.(types.String).String())
	case types.BlobKind:
		panic("not yet implemented")
	case types.ValueKind:
		// The value is always tagged
		runtimeType := v.TypeRef()
		w.writeTypeRef(runtimeType)
		w.writeValue(runtimeType, v)
	case types.ListKind:
		w.writeList(t, v.(types.List))
	case types.MapKind:
		w.writeMap(t, v.(types.Map))
	case types.RefKind:
		panic("not yet implemented")
	case types.SetKind:
		w.writeSet(t, v.(types.Set))
	case types.EnumKind:
		w.writeEnum(t, v.(types.UInt32))
	case types.StructKind:
		w.writeStruct(t, v.(types.Map))
	case types.TypeRefKind:
		panic("not yet implemented")
	}
}

func (w *jsonArrayWriter) writeList(t types.TypeRef, l types.List) {
	desc := t.Desc.(types.CompoundDesc)
	elemType := desc.ElemTypes[0]
	l.IterAll(func(v types.Value) {
		w.writeValue(elemType, v)
	})
}

func (w *jsonArrayWriter) writeSet(t types.TypeRef, l types.Set) {
	desc := t.Desc.(types.CompoundDesc)
	elemType := desc.ElemTypes[0]
	l.IterAll(func(v types.Value) {
		w.writeValue(elemType, v)
	})
}

func (w *jsonArrayWriter) writeMap(t types.TypeRef, l types.Map) {
	desc := t.Desc.(types.CompoundDesc)
	keyType := desc.ElemTypes[0]
	valueType := desc.ElemTypes[1]
	l.IterAll(func(k, v types.Value) {
		w.writeValue(keyType, k)
		w.writeValue(valueType, v)
	})
}

func (w *jsonArrayWriter) writeStruct(t types.TypeRef, m types.Map) {
	desc := t.Desc.(types.StructDesc)
	for _, f := range desc.Fields {
		v, ok := m.MaybeGet(types.NewString(f.Name))
		if f.Optional {
			if ok {
				w.write(uint32(1))
				w.writeValue(f.T, v)
			} else {
				w.write(uint32(0))
				w.write(uint32(0))
			}
		} else {
			d.Chk.True(ok)
			w.writeValue(f.T, v)
		}
	}
	if len(desc.Union) > 0 {
		i := uint32(m.Get(types.NewString("$unionIndex")).(types.UInt32))
		v := m.Get(types.NewString("$unionValue"))
		w.write(i)
		w.writeValue(desc.Union[i].T, v)
	}
}

func (w *jsonArrayWriter) writeEnum(t types.TypeRef, v types.UInt32) {
	w.write(uint32(v))
}
