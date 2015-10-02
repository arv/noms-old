// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_ref_CachedRef = __testPackageInFile_ref_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_ref_Ref() ref.Ref {
	p := types.PackageDef{
		NamedTypes: types.MapOfStringToTypeRefDef{

			"StructWithRef": types.MakeStructTypeRef("StructWithRef",
				[]types.Field{
					types.Field{"r", types.MakeCompoundTypeRef("", types.RefKind, types.MakeCompoundTypeRef("", types.SetKind, types.MakePrimitiveTypeRef(types.Float32Kind))), false},
				},
				types.Choices{},
			),
		},
	}.New()
	return types.RegisterPackage(&p)
}

// RefOfListOfString

type RefOfListOfString struct {
	r ref.Ref
}

func NewRefOfListOfString(r ref.Ref) RefOfListOfString {
	return RefOfListOfString{r}
}

func (r RefOfListOfString) Ref() ref.Ref {
	return r.r
}

func (r RefOfListOfString) Equals(other RefOfListOfString) bool {
	return r.Ref() == other.Ref()
}

// A Noms Value that describes RefOfListOfString.
var __typeRefForRefOfListOfString = types.MakeCompoundTypeRef("", types.RefKind, types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.StringKind)))

func (m RefOfListOfString) TypeRef() types.TypeRef {
	return __typeRefForRefOfListOfString
}

func (r RefOfListOfString) NomsValue() types.Value {
	return types.Ref{R: r.r}
}

func RefOfListOfStringFromVal(p types.Value) RefOfListOfString {
	return RefOfListOfString{p.(types.Ref).Ref()}
}

func (r RefOfListOfString) GetValue(cs chunks.ChunkSource) ListOfString {
	return ListOfStringFromVal(types.ReadValue(r.r, cs))
}

func (r RefOfListOfString) SetValue(val ListOfString, cs chunks.ChunkSink) RefOfListOfString {
	ref := types.WriteValue(val.NomsValue(), cs)
	return RefOfListOfString{ref}
}

// ListOfString

type ListOfString struct {
	l types.List
}

func NewListOfString() ListOfString {
	return ListOfString{types.NewList()}
}

type ListOfStringDef []string

func (def ListOfStringDef) New() ListOfString {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = types.NewString(d)
	}
	return ListOfString{types.NewList(l...)}
}

func (l ListOfString) Def() ListOfStringDef {
	d := make([]string, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = l.l.Get(i).(types.String).String()
	}
	return d
}

func ListOfStringFromVal(val types.Value) ListOfString {
	// TODO: Validate here
	return ListOfString{val.(types.List)}
}

func (l ListOfString) NomsValue() types.Value {
	return l.l
}

func (l ListOfString) Equals(p ListOfString) bool {
	return l.l.Equals(p.l)
}

func (l ListOfString) Ref() ref.Ref {
	return l.l.Ref()
}

// A Noms Value that describes ListOfString.
var __typeRefForListOfString = types.MakeCompoundTypeRef("", types.ListKind, types.MakePrimitiveTypeRef(types.StringKind))

func (m ListOfString) TypeRef() types.TypeRef {
	return __typeRefForListOfString
}

func (l ListOfString) Len() uint64 {
	return l.l.Len()
}

func (l ListOfString) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfString) Get(i uint64) string {
	return l.l.Get(i).(types.String).String()
}

func (l ListOfString) Slice(idx uint64, end uint64) ListOfString {
	return ListOfString{l.l.Slice(idx, end)}
}

func (l ListOfString) Set(i uint64, val string) ListOfString {
	return ListOfString{l.l.Set(i, types.NewString(val))}
}

func (l ListOfString) Append(v ...string) ListOfString {
	return ListOfString{l.l.Append(l.fromElemSlice(v)...)}
}

func (l ListOfString) Insert(idx uint64, v ...string) ListOfString {
	return ListOfString{l.l.Insert(idx, l.fromElemSlice(v)...)}
}

func (l ListOfString) Remove(idx uint64, end uint64) ListOfString {
	return ListOfString{l.l.Remove(idx, end)}
}

func (l ListOfString) RemoveAt(idx uint64) ListOfString {
	return ListOfString{(l.l.RemoveAt(idx))}
}

func (l ListOfString) fromElemSlice(p []string) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = types.NewString(v)
	}
	return r
}

type ListOfStringIterCallback func(v string, i uint64) (stop bool)

func (l ListOfString) Iter(cb ListOfStringIterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(v.(types.String).String(), i)
	})
}

type ListOfStringIterAllCallback func(v string, i uint64)

func (l ListOfString) IterAll(cb ListOfStringIterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(v.(types.String).String(), i)
	})
}

type ListOfStringFilterCallback func(v string, i uint64) (keep bool)

func (l ListOfString) Filter(cb ListOfStringFilterCallback) ListOfString {
	nl := NewListOfString()
	l.IterAll(func(v string, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}

// ListOfRefOfFloat32

type ListOfRefOfFloat32 struct {
	l types.List
}

func NewListOfRefOfFloat32() ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{types.NewList()}
}

type ListOfRefOfFloat32Def []ref.Ref

func (def ListOfRefOfFloat32Def) New() ListOfRefOfFloat32 {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = types.Ref{R: d}
	}
	return ListOfRefOfFloat32{types.NewList(l...)}
}

func (l ListOfRefOfFloat32) Def() ListOfRefOfFloat32Def {
	d := make([]ref.Ref, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = l.l.Get(i).Ref()
	}
	return d
}

func ListOfRefOfFloat32FromVal(val types.Value) ListOfRefOfFloat32 {
	// TODO: Validate here
	return ListOfRefOfFloat32{val.(types.List)}
}

func (l ListOfRefOfFloat32) NomsValue() types.Value {
	return l.l
}

func (l ListOfRefOfFloat32) Equals(p ListOfRefOfFloat32) bool {
	return l.l.Equals(p.l)
}

func (l ListOfRefOfFloat32) Ref() ref.Ref {
	return l.l.Ref()
}

// A Noms Value that describes ListOfRefOfFloat32.
var __typeRefForListOfRefOfFloat32 = types.MakeCompoundTypeRef("", types.ListKind, types.MakeCompoundTypeRef("", types.RefKind, types.MakePrimitiveTypeRef(types.Float32Kind)))

func (m ListOfRefOfFloat32) TypeRef() types.TypeRef {
	return __typeRefForListOfRefOfFloat32
}

func (l ListOfRefOfFloat32) Len() uint64 {
	return l.l.Len()
}

func (l ListOfRefOfFloat32) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfRefOfFloat32) Get(i uint64) RefOfFloat32 {
	return RefOfFloat32FromVal(l.l.Get(i))
}

func (l ListOfRefOfFloat32) Slice(idx uint64, end uint64) ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{l.l.Slice(idx, end)}
}

func (l ListOfRefOfFloat32) Set(i uint64, val RefOfFloat32) ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{l.l.Set(i, val.NomsValue())}
}

func (l ListOfRefOfFloat32) Append(v ...RefOfFloat32) ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{l.l.Append(l.fromElemSlice(v)...)}
}

func (l ListOfRefOfFloat32) Insert(idx uint64, v ...RefOfFloat32) ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{l.l.Insert(idx, l.fromElemSlice(v)...)}
}

func (l ListOfRefOfFloat32) Remove(idx uint64, end uint64) ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{l.l.Remove(idx, end)}
}

func (l ListOfRefOfFloat32) RemoveAt(idx uint64) ListOfRefOfFloat32 {
	return ListOfRefOfFloat32{(l.l.RemoveAt(idx))}
}

func (l ListOfRefOfFloat32) fromElemSlice(p []RefOfFloat32) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v.NomsValue()
	}
	return r
}

type ListOfRefOfFloat32IterCallback func(v RefOfFloat32, i uint64) (stop bool)

func (l ListOfRefOfFloat32) Iter(cb ListOfRefOfFloat32IterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(RefOfFloat32FromVal(v), i)
	})
}

type ListOfRefOfFloat32IterAllCallback func(v RefOfFloat32, i uint64)

func (l ListOfRefOfFloat32) IterAll(cb ListOfRefOfFloat32IterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(RefOfFloat32FromVal(v), i)
	})
}

type ListOfRefOfFloat32FilterCallback func(v RefOfFloat32, i uint64) (keep bool)

func (l ListOfRefOfFloat32) Filter(cb ListOfRefOfFloat32FilterCallback) ListOfRefOfFloat32 {
	nl := NewListOfRefOfFloat32()
	l.IterAll(func(v RefOfFloat32, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}

// RefOfFloat32

type RefOfFloat32 struct {
	r ref.Ref
}

func NewRefOfFloat32(r ref.Ref) RefOfFloat32 {
	return RefOfFloat32{r}
}

func (r RefOfFloat32) Ref() ref.Ref {
	return r.r
}

func (r RefOfFloat32) Equals(other RefOfFloat32) bool {
	return r.Ref() == other.Ref()
}

// A Noms Value that describes RefOfFloat32.
var __typeRefForRefOfFloat32 = types.MakeCompoundTypeRef("", types.RefKind, types.MakePrimitiveTypeRef(types.Float32Kind))

func (m RefOfFloat32) TypeRef() types.TypeRef {
	return __typeRefForRefOfFloat32
}

func (r RefOfFloat32) NomsValue() types.Value {
	return types.Ref{R: r.r}
}

func RefOfFloat32FromVal(p types.Value) RefOfFloat32 {
	return RefOfFloat32{p.(types.Ref).Ref()}
}

func (r RefOfFloat32) GetValue(cs chunks.ChunkSource) float32 {
	return float32(types.ReadValue(r.r, cs).(types.Float32))
}

func (r RefOfFloat32) SetValue(val float32, cs chunks.ChunkSink) RefOfFloat32 {
	ref := types.WriteValue(types.Float32(val), cs)
	return RefOfFloat32{ref}
}

// StructWithRef

type StructWithRef struct {
	m types.Map
}

func NewStructWithRef() StructWithRef {
	return StructWithRef{types.NewMap(
		types.NewString("$name"), types.NewString("StructWithRef"),
		types.NewString("$type"), types.MakeTypeRef("StructWithRef", __testPackageInFile_ref_CachedRef),
		types.NewString("r"), types.Ref{R: ref.Ref{}},
	)}
}

type StructWithRefDef struct {
	R ref.Ref
}

func (def StructWithRefDef) New() StructWithRef {
	return StructWithRef{
		types.NewMap(
			types.NewString("$name"), types.NewString("StructWithRef"),
			types.NewString("$type"), types.MakeTypeRef("StructWithRef", __testPackageInFile_ref_CachedRef),
			types.NewString("r"), types.Ref{R: def.R},
		)}
}

func (s StructWithRef) Def() (d StructWithRefDef) {
	d.R = s.m.Get(types.NewString("r")).Ref()
	return
}

func (m StructWithRef) TypeRef() types.TypeRef {
	return types.MakeTypeRef("StructWithRef", __testPackageInFile_ref_CachedRef)
}

func StructWithRefFromVal(val types.Value) StructWithRef {
	// TODO: Validate here
	return StructWithRef{val.(types.Map)}
}

func (s StructWithRef) NomsValue() types.Value {
	return s.m
}

func (s StructWithRef) Equals(other StructWithRef) bool {
	return s.m.Equals(other.m)
}

func (s StructWithRef) Ref() ref.Ref {
	return s.m.Ref()
}

func (s StructWithRef) R() RefOfSetOfFloat32 {
	return RefOfSetOfFloat32FromVal(s.m.Get(types.NewString("r")))
}

func (s StructWithRef) SetR(val RefOfSetOfFloat32) StructWithRef {
	return StructWithRef{s.m.Set(types.NewString("r"), val.NomsValue())}
}

// RefOfSetOfFloat32

type RefOfSetOfFloat32 struct {
	r ref.Ref
}

func NewRefOfSetOfFloat32(r ref.Ref) RefOfSetOfFloat32 {
	return RefOfSetOfFloat32{r}
}

func (r RefOfSetOfFloat32) Ref() ref.Ref {
	return r.r
}

func (r RefOfSetOfFloat32) Equals(other RefOfSetOfFloat32) bool {
	return r.Ref() == other.Ref()
}

// A Noms Value that describes RefOfSetOfFloat32.
var __typeRefForRefOfSetOfFloat32 = types.MakeCompoundTypeRef("", types.RefKind, types.MakeCompoundTypeRef("", types.SetKind, types.MakePrimitiveTypeRef(types.Float32Kind)))

func (m RefOfSetOfFloat32) TypeRef() types.TypeRef {
	return __typeRefForRefOfSetOfFloat32
}

func (r RefOfSetOfFloat32) NomsValue() types.Value {
	return types.Ref{R: r.r}
}

func RefOfSetOfFloat32FromVal(p types.Value) RefOfSetOfFloat32 {
	return RefOfSetOfFloat32{p.(types.Ref).Ref()}
}

func (r RefOfSetOfFloat32) GetValue(cs chunks.ChunkSource) SetOfFloat32 {
	return SetOfFloat32FromVal(types.ReadValue(r.r, cs))
}

func (r RefOfSetOfFloat32) SetValue(val SetOfFloat32, cs chunks.ChunkSink) RefOfSetOfFloat32 {
	ref := types.WriteValue(val.NomsValue(), cs)
	return RefOfSetOfFloat32{ref}
}

// SetOfFloat32

type SetOfFloat32 struct {
	s types.Set
}

func NewSetOfFloat32() SetOfFloat32 {
	return SetOfFloat32{types.NewSet()}
}

type SetOfFloat32Def map[float32]bool

func (def SetOfFloat32Def) New() SetOfFloat32 {
	l := make([]types.Value, len(def))
	i := 0
	for d, _ := range def {
		l[i] = types.Float32(d)
		i++
	}
	return SetOfFloat32{types.NewSet(l...)}
}

func (s SetOfFloat32) Def() SetOfFloat32Def {
	def := make(map[float32]bool, s.Len())
	s.s.Iter(func(v types.Value) bool {
		def[float32(v.(types.Float32))] = true
		return false
	})
	return def
}

func SetOfFloat32FromVal(p types.Value) SetOfFloat32 {
	return SetOfFloat32{p.(types.Set)}
}

func (s SetOfFloat32) NomsValue() types.Value {
	return s.s
}

func (s SetOfFloat32) Equals(p SetOfFloat32) bool {
	return s.s.Equals(p.s)
}

func (s SetOfFloat32) Ref() ref.Ref {
	return s.s.Ref()
}

// A Noms Value that describes SetOfFloat32.
var __typeRefForSetOfFloat32 = types.MakeCompoundTypeRef("", types.SetKind, types.MakePrimitiveTypeRef(types.Float32Kind))

func (m SetOfFloat32) TypeRef() types.TypeRef {
	return __typeRefForSetOfFloat32
}

func (s SetOfFloat32) Empty() bool {
	return s.s.Empty()
}

func (s SetOfFloat32) Len() uint64 {
	return s.s.Len()
}

func (s SetOfFloat32) Has(p float32) bool {
	return s.s.Has(types.Float32(p))
}

type SetOfFloat32IterCallback func(p float32) (stop bool)

func (s SetOfFloat32) Iter(cb SetOfFloat32IterCallback) {
	s.s.Iter(func(v types.Value) bool {
		return cb(float32(v.(types.Float32)))
	})
}

type SetOfFloat32IterAllCallback func(p float32)

func (s SetOfFloat32) IterAll(cb SetOfFloat32IterAllCallback) {
	s.s.IterAll(func(v types.Value) {
		cb(float32(v.(types.Float32)))
	})
}

type SetOfFloat32FilterCallback func(p float32) (keep bool)

func (s SetOfFloat32) Filter(cb SetOfFloat32FilterCallback) SetOfFloat32 {
	ns := NewSetOfFloat32()
	s.IterAll(func(v float32) {
		if cb(v) {
			ns = ns.Insert(v)
		}
	})
	return ns
}

func (s SetOfFloat32) Insert(p ...float32) SetOfFloat32 {
	return SetOfFloat32{s.s.Insert(s.fromElemSlice(p)...)}
}

func (s SetOfFloat32) Remove(p ...float32) SetOfFloat32 {
	return SetOfFloat32{s.s.Remove(s.fromElemSlice(p)...)}
}

func (s SetOfFloat32) Union(others ...SetOfFloat32) SetOfFloat32 {
	return SetOfFloat32{s.s.Union(s.fromStructSlice(others)...)}
}

func (s SetOfFloat32) Subtract(others ...SetOfFloat32) SetOfFloat32 {
	return SetOfFloat32{s.s.Subtract(s.fromStructSlice(others)...)}
}

func (s SetOfFloat32) Any() float32 {
	return float32(s.s.Any().(types.Float32))
}

func (s SetOfFloat32) fromStructSlice(p []SetOfFloat32) []types.Set {
	r := make([]types.Set, len(p))
	for i, v := range p {
		r[i] = v.s
	}
	return r
}

func (s SetOfFloat32) fromElemSlice(p []float32) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = types.Float32(v)
	}
	return r
}
