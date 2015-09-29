package encode

import (
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/types"
)

type jsonArrayWriter []interface{}

func newJsonArrayWriter() *jsonArrayWriter {
	return &jsonArrayWriter{}
}

func (w *jsonArrayWriter) write(v interface{}) {
	*w = append(*w, v)
}

func (w *jsonArrayWriter) toSlice() []interface{} {
	return *w
}

func (w *jsonArrayWriter) writeRef(r types.Ref) {
	w.write(r.Ref().String())
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

func (w *jsonArrayWriter) writeValue(t types.TypeRef, v types.Value) {
	switch t.Kind() {
	case types.BoolKind:
		w.write(bool(v.(types.Bool)))
	case types.UInt8Kind:
		w.write(uint8(v.(types.UInt8)))
	case types.UInt16Kind:
		w.write(uint16(v.(types.UInt16)))
	case types.UInt32Kind:
		w.write(uint32(v.(types.UInt32)))
	case types.UInt64Kind:
		w.write(uint64(v.(types.UInt64)))
	case types.Int8Kind:
		w.write(int8(v.(types.Int8)))
	case types.Int16Kind:
		w.write(int16(v.(types.Int16)))
	case types.Int32Kind:
		w.write(int32(v.(types.Int32)))
	case types.Int64Kind:
		w.write(int64(v.(types.Int64)))
	case types.Float32Kind:
		w.write(float32(v.(types.Float32)))
	case types.Float64Kind:
		w.write(float64(v.(types.Float64)))
	case types.StringKind:
		w.write(v.(types.String).String())
	case types.BlobKind:
		panic("not yet implemented")
	case types.ValueKind:
		// The value is always tagged
		w.writeTypeRef(t)
		panic("not yet implemented")
	case types.ListKind:
		w2 := newJsonArrayWriter()
		w2.writeList(t, v.(types.List))
		w.write(w2.toSlice())
	case types.MapKind:
		w2 := newJsonArrayWriter()
		w2.writeMap(t, v.(types.Map))
		w.write(w2.toSlice())
	case types.RefKind:
		panic("not yet implemented")
	case types.SetKind:
		w2 := newJsonArrayWriter()
		w2.writeSet(t, v.(types.Set))
		w.write(w2.toSlice())
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
