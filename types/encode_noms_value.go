package types

import (
	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/ref"
)

// typedValue implements enc.typedValue which is used to tag the value for now so that we can trigger a different encoding strategy.
type typedValue struct {
	v interface{}
}

func (tv typedValue) TypedValue() interface{} {
	return tv.v
}

func encNomsValue(v NomsValue, cs chunks.ChunkSink) interface{} {
	w := newJsonArrayWriter()
	t := v.TypeRef()
	w.writeTypeRef(t)
	w.writeTopLevelValue(t, v.NomsValue())
	return typedValue{w.toArray()}
}

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

func (w *jsonArrayWriter) writeTypeRef(t TypeRef) {
	// TODO: Resolve if needed
	k := t.Kind()
	w.write(k)
	switch k {
	case EnumKind, StructKind:
		w.writeRef(t.PackageRef())
		// TODO: Should be ordinal instead of name.
		w.write(t.Name())
	case ListKind, MapKind, RefKind, SetKind:
		for _, elemType := range t.Desc.(CompoundDesc).ElemTypes {
			w.writeTypeRef(elemType)
		}
	}
}

func (w *jsonArrayWriter) writeNomsValue(nv NomsValue) {
	v := nv.NomsValue()
	t := nv.TypeRef()
	w.writeTypeRef(t)
	w.writeTopLevelValue(t, v)
}

func (w *jsonArrayWriter) writeValue(t TypeRef, v Value) {
	switch t.Kind() {
	case ListKind, MapKind, SetKind:
		w2 := newJsonArrayWriter()
		w2.writeTopLevelValue(t, v)
		w.write(w2.toArray())
	default:
		w.writeTopLevelValue(t, v)
	}
}

func (w *jsonArrayWriter) writeTopLevelValue(t TypeRef, v Value) {
	switch t.Kind() {
	case BoolKind, Float32Kind, Float64Kind, Int16Kind, Int32Kind, Int64Kind, Int8Kind, UInt16Kind, UInt32Kind, UInt64Kind, UInt8Kind:
		w.write(v.(primitive).ToPrimitive())
	case StringKind:
		w.write(v.(String).String())
	case BlobKind:
		panic("not yet implemented")
	case ValueKind:
		// The value is always tagged
		runtimeType := v.TypeRef()
		w.writeTypeRef(runtimeType)
		w.writeValue(runtimeType, v)
	case ListKind:
		w.writeList(t, v.(List))
	case MapKind:
		w.writeMap(t, v.(Map))
	case RefKind:
		panic("not yet implemented")
	case SetKind:
		w.writeSet(t, v.(Set))
	case EnumKind:
		w.writeEnum(t, v.(UInt32))
	case StructKind:
		w.writeStruct(t, v.(Map))
	case TypeRefKind:
		panic("not yet implemented")
	}
}

func (w *jsonArrayWriter) writeList(t TypeRef, l List) {
	desc := t.Desc.(CompoundDesc)
	elemType := desc.ElemTypes[0]
	l.IterAll(func(v Value, i uint64) {
		w.writeValue(elemType, v)
	})
}

func (w *jsonArrayWriter) writeSet(t TypeRef, s Set) {
	desc := t.Desc.(CompoundDesc)
	elemType := desc.ElemTypes[0]
	s.IterAll(func(v Value) {
		w.writeValue(elemType, v)
	})
}

func (w *jsonArrayWriter) writeMap(t TypeRef, m Map) {
	desc := t.Desc.(CompoundDesc)
	keyType := desc.ElemTypes[0]
	valueType := desc.ElemTypes[1]
	m.IterAll(func(k, v Value) {
		w.writeValue(keyType, k)
		w.writeValue(valueType, v)
	})
}

func (w *jsonArrayWriter) writeStruct(t TypeRef, m Map) {
	desc := t.Desc.(StructDesc)
	for _, f := range desc.Fields {
		v, ok := m.MaybeGet(NewString(f.Name))
		if f.Optional {
			if ok {
				w.write(true)
				w.writeValue(f.T, v)
			} else {
				w.write(false)
			}
		} else {
			d.Chk.True(ok)
			w.writeValue(f.T, v)
		}
	}
	if len(desc.Union) > 0 {
		i := uint32(m.Get(NewString("$unionIndex")).(UInt32))
		v := m.Get(NewString("$unionValue"))
		w.write(i)
		w.writeValue(desc.Union[i].T, v)
	}
}

func (w *jsonArrayWriter) writeEnum(t TypeRef, v UInt32) {
	w.write(uint32(v))
}
