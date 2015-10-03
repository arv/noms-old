// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_struct_with_list_CachedRef = __testPackageInFile_struct_with_list_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_struct_with_list_Ref() ref.Ref {
	p := types.PackageDef{
		NamedTypes: types.MapOfStringToTypeRefDef{

			"StructWithList": types.MakeStructTypeRef("StructWithList",
				[]types.Field{
					types.Field{"l", types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.UInt8Kind)), false},
					types.Field{"b", types.MakePrimitiveTypeRef(types.BoolKind), false},
					types.Field{"s", types.MakePrimitiveTypeRef(types.StringKind), false},
					types.Field{"i", types.MakePrimitiveTypeRef(types.Int64Kind), false},
				},
				types.Choices{},
			),
		},
	}.New()
	return types.RegisterPackage(&p)
}

// StructWithList

type StructWithList struct {
	m types.Map
}

func NewStructWithList() StructWithList {
	return StructWithList{types.NewMap(
		types.NewString("$name"), types.NewString("StructWithList"),
		types.NewString("$type"), types.MakeTypeRef("StructWithList", __testPackageInFile_struct_with_list_CachedRef),
		types.NewString("l"), types.NewList(),
		types.NewString("b"), types.Bool(false),
		types.NewString("s"), types.NewString(""),
		types.NewString("i"), types.Int64(0),
	)}
}

type StructWithListDef struct {
	L ListOfUInt8Def
	B bool
	S string
	I int64
}

func (def StructWithListDef) New() StructWithList {
	return StructWithList{
		types.NewMap(
			types.NewString("$name"), types.NewString("StructWithList"),
			types.NewString("$type"), types.MakeTypeRef("StructWithList", __testPackageInFile_struct_with_list_CachedRef),
			types.NewString("l"), def.L.New().NomsValue(),
			types.NewString("b"), types.Bool(def.B),
			types.NewString("s"), types.NewString(def.S),
			types.NewString("i"), types.Int64(def.I),
		)}
}

func (s StructWithList) Def() (d StructWithListDef) {
	d.L = ListOfUInt8FromVal(s.m.Get(types.NewString("l"))).Def()
	d.B = bool(s.m.Get(types.NewString("b")).(types.Bool))
	d.S = s.m.Get(types.NewString("s")).(types.String).String()
	d.I = int64(s.m.Get(types.NewString("i")).(types.Int64))
	return
}

var __typeRefForStructWithList = types.MakeTypeRef("StructWithList", __testPackageInFile_struct_with_list_CachedRef)

func (m StructWithList) TypeRef() types.TypeRef {
	return __typeRefForStructWithList
}

func StructWithListFromVal(val types.Value) StructWithList {
	// TODO: Validate here
	return StructWithList{val.(types.Map)}
}

func (s StructWithList) NomsValue() types.Value {
	return s.m
}

func (s StructWithList) Equals(other StructWithList) bool {
	return s.m.Equals(other.m)
}

func (s StructWithList) Ref() ref.Ref {
	return s.m.Ref()
}

func (s StructWithList) L() ListOfUInt8 {
	return ListOfUInt8FromVal(s.m.Get(types.NewString("l")))
}

func (s StructWithList) SetL(val ListOfUInt8) StructWithList {
	return StructWithList{s.m.Set(types.NewString("l"), val.NomsValue())}
}

func (s StructWithList) B() bool {
	return bool(s.m.Get(types.NewString("b")).(types.Bool))
}

func (s StructWithList) SetB(val bool) StructWithList {
	return StructWithList{s.m.Set(types.NewString("b"), types.Bool(val))}
}

func (s StructWithList) S() string {
	return s.m.Get(types.NewString("s")).(types.String).String()
}

func (s StructWithList) SetS(val string) StructWithList {
	return StructWithList{s.m.Set(types.NewString("s"), types.NewString(val))}
}

func (s StructWithList) I() int64 {
	return int64(s.m.Get(types.NewString("i")).(types.Int64))
}

func (s StructWithList) SetI(val int64) StructWithList {
	return StructWithList{s.m.Set(types.NewString("i"), types.Int64(val))}
}

// ListOfUInt8

type ListOfUInt8 struct {
	l types.List
}

func NewListOfUInt8() ListOfUInt8 {
	return ListOfUInt8{types.NewList()}
}

type ListOfUInt8Def []uint8

func (def ListOfUInt8Def) New() ListOfUInt8 {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = types.UInt8(d)
	}
	return ListOfUInt8{types.NewList(l...)}
}

func (l ListOfUInt8) Def() ListOfUInt8Def {
	d := make([]uint8, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = uint8(l.l.Get(i).(types.UInt8))
	}
	return d
}

func ListOfUInt8FromVal(val types.Value) ListOfUInt8 {
	// TODO: Validate here
	return ListOfUInt8{val.(types.List)}
}

func (l ListOfUInt8) NomsValue() types.Value {
	return l.l
}

func (l ListOfUInt8) Equals(p ListOfUInt8) bool {
	return l.l.Equals(p.l)
}

func (l ListOfUInt8) Ref() ref.Ref {
	return l.l.Ref()
}

// A Noms Value that describes ListOfUInt8.
var __typeRefForListOfUInt8 = types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.UInt8Kind))

func (m ListOfUInt8) TypeRef() types.TypeRef {
	return __typeRefForListOfUInt8
}

func (l ListOfUInt8) Len() uint64 {
	return l.l.Len()
}

func (l ListOfUInt8) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfUInt8) Get(i uint64) uint8 {
	return uint8(l.l.Get(i).(types.UInt8))
}

func (l ListOfUInt8) Slice(idx uint64, end uint64) ListOfUInt8 {
	return ListOfUInt8{l.l.Slice(idx, end)}
}

func (l ListOfUInt8) Set(i uint64, val uint8) ListOfUInt8 {
	return ListOfUInt8{l.l.Set(i, types.UInt8(val))}
}

func (l ListOfUInt8) Append(v ...uint8) ListOfUInt8 {
	return ListOfUInt8{l.l.Append(l.fromElemSlice(v)...)}
}

func (l ListOfUInt8) Insert(idx uint64, v ...uint8) ListOfUInt8 {
	return ListOfUInt8{l.l.Insert(idx, l.fromElemSlice(v)...)}
}

func (l ListOfUInt8) Remove(idx uint64, end uint64) ListOfUInt8 {
	return ListOfUInt8{l.l.Remove(idx, end)}
}

func (l ListOfUInt8) RemoveAt(idx uint64) ListOfUInt8 {
	return ListOfUInt8{(l.l.RemoveAt(idx))}
}

func (l ListOfUInt8) fromElemSlice(p []uint8) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = types.UInt8(v)
	}
	return r
}

type ListOfUInt8IterCallback func(v uint8, i uint64) (stop bool)

func (l ListOfUInt8) Iter(cb ListOfUInt8IterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(uint8(v.(types.UInt8)), i)
	})
}

type ListOfUInt8IterAllCallback func(v uint8, i uint64)

func (l ListOfUInt8) IterAll(cb ListOfUInt8IterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(uint8(v.(types.UInt8)), i)
	})
}

type ListOfUInt8FilterCallback func(v uint8, i uint64) (keep bool)

func (l ListOfUInt8) Filter(cb ListOfUInt8FilterCallback) ListOfUInt8 {
	nl := NewListOfUInt8()
	l.IterAll(func(v uint8, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}
