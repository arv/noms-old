// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_struct_CachedRef = __testPackageInFile_struct_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_struct_Ref() ref.Ref {
	p := types.PackageDef{
		NamedTypes: types.MapOfStringToTypeRefDef{

			"Struct": types.MakeStructTypeRef("Struct",
				[]types.Field{
					types.Field{"s", types.MakePrimitiveTypeRef(types.StringKind), false},
					types.Field{"b", types.MakePrimitiveTypeRef(types.BoolKind), false},
				},
				types.Choices{},
			),
		},
	}.New()
	return types.RegisterPackage(&p)
}

// Struct

type Struct struct {
	m types.Map
}

func NewStruct() Struct {
	return Struct{types.NewMap(
		types.NewString("$name"), types.NewString("Struct"),
		types.NewString("$type"), types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef),
		types.NewString("s"), types.NewString(""),
		types.NewString("b"), types.Bool(false),
	)}
}

type StructDef struct {
	S string
	B bool
}

func (def StructDef) New() Struct {
	return Struct{
		types.NewMap(
			types.NewString("$name"), types.NewString("Struct"),
			types.NewString("$type"), types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef),
			types.NewString("s"), types.NewString(def.S),
			types.NewString("b"), types.Bool(def.B),
		)}
}

func (s Struct) Def() (d StructDef) {
	d.S = s.m.Get(types.NewString("s")).(types.String).String()
	d.B = bool(s.m.Get(types.NewString("b")).(types.Bool))
	return
}

func (m Struct) TypeRef() types.TypeRef {
	return types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef)
}

func StructFromVal(val types.Value) Struct {
	// TODO: Validate here
	return Struct{val.(types.Map)}
}

func (s Struct) NomsValue() types.Value {
	return s.m
}

func (s Struct) Equals(other Struct) bool {
	return s.m.Equals(other.m)
}

func (s Struct) Ref() ref.Ref {
	return s.m.Ref()
}

func (s Struct) S() string {
	return s.m.Get(types.NewString("s")).(types.String).String()
}

func (s Struct) SetS(val string) Struct {
	return Struct{s.m.Set(types.NewString("s"), types.NewString(val))}
}

func (s Struct) B() bool {
	return bool(s.m.Get(types.NewString("b")).(types.Bool))
}

func (s Struct) SetB(val bool) Struct {
	return Struct{s.m.Set(types.NewString("b"), types.Bool(val))}
}
