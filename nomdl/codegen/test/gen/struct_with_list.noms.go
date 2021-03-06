// This file was generated by nomdl/codegen.

package gen

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __genPackageInFile_struct_with_list_CachedRef ref.Ref

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func init() {
	p := types.NewPackage([]types.Type{
		types.MakeStructType("StructWithList",
			[]types.Field{
				types.Field{"l", types.MakeCompoundType(types.ListKind, types.MakePrimitiveType(types.Uint8Kind)), false},
				types.Field{"b", types.MakePrimitiveType(types.BoolKind), false},
				types.Field{"s", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"i", types.MakePrimitiveType(types.Int64Kind), false},
			},
			types.Choices{},
		),
	}, []ref.Ref{})
	__genPackageInFile_struct_with_list_CachedRef = types.RegisterPackage(&p)
}

// StructWithList

type StructWithList struct {
	_l ListOfUint8
	_b bool
	_s string
	_i int64

	ref *ref.Ref
}

func NewStructWithList() StructWithList {
	return StructWithList{
		_l: NewListOfUint8(),
		_b: false,
		_s: "",
		_i: int64(0),

		ref: &ref.Ref{},
	}
}

type StructWithListDef struct {
	L ListOfUint8Def
	B bool
	S string
	I int64
}

func (def StructWithListDef) New() StructWithList {
	return StructWithList{
		_l:  def.L.New(),
		_b:  def.B,
		_s:  def.S,
		_i:  def.I,
		ref: &ref.Ref{},
	}
}

func (s StructWithList) Def() (d StructWithListDef) {
	d.L = s._l.Def()
	d.B = s._b
	d.S = s._s
	d.I = s._i
	return
}

var __typeForStructWithList types.Type

func (m StructWithList) Type() types.Type {
	return __typeForStructWithList
}

func init() {
	__typeForStructWithList = types.MakeType(__genPackageInFile_struct_with_list_CachedRef, 0)
	types.RegisterStruct(__typeForStructWithList, builderForStructWithList, readerForStructWithList)
}

func builderForStructWithList(values []types.Value) types.Value {
	i := 0
	s := StructWithList{ref: &ref.Ref{}}
	s._l = values[i].(ListOfUint8)
	i++
	s._b = bool(values[i].(types.Bool))
	i++
	s._s = values[i].(types.String).String()
	i++
	s._i = int64(values[i].(types.Int64))
	i++
	return s
}

func readerForStructWithList(v types.Value) []types.Value {
	values := []types.Value{}
	s := v.(StructWithList)
	values = append(values, s._l)
	values = append(values, types.Bool(s._b))
	values = append(values, types.NewString(s._s))
	values = append(values, types.Int64(s._i))
	return values
}

func (s StructWithList) Equals(other types.Value) bool {
	return other != nil && __typeForStructWithList.Equals(other.Type()) && s.Ref() == other.Ref()
}

func (s StructWithList) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s StructWithList) Chunks() (chunks []types.RefBase) {
	chunks = append(chunks, __typeForStructWithList.Chunks()...)
	chunks = append(chunks, s._l.Chunks()...)
	return
}

func (s StructWithList) ChildValues() (ret []types.Value) {
	ret = append(ret, s._l)
	ret = append(ret, types.Bool(s._b))
	ret = append(ret, types.NewString(s._s))
	ret = append(ret, types.Int64(s._i))
	return
}

func (s StructWithList) L() ListOfUint8 {
	return s._l
}

func (s StructWithList) SetL(val ListOfUint8) StructWithList {
	s._l = val
	s.ref = &ref.Ref{}
	return s
}

func (s StructWithList) B() bool {
	return s._b
}

func (s StructWithList) SetB(val bool) StructWithList {
	s._b = val
	s.ref = &ref.Ref{}
	return s
}

func (s StructWithList) S() string {
	return s._s
}

func (s StructWithList) SetS(val string) StructWithList {
	s._s = val
	s.ref = &ref.Ref{}
	return s
}

func (s StructWithList) I() int64 {
	return s._i
}

func (s StructWithList) SetI(val int64) StructWithList {
	s._i = val
	s.ref = &ref.Ref{}
	return s
}
