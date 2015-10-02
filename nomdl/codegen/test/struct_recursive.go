// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_struct_recursive_CachedRef = __testPackageInFile_struct_recursive_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_struct_recursive_Ref() ref.Ref {
	p := types.PackageDef{
		NamedTypes: types.MapOfStringToTypeRefDef{

			"Tree": types.MakeStructTypeRef("Tree",
				[]types.Field{
					types.Field{"children", types.MakeCompoundTypeRef("", types.ListKind, types.MakeTypeRef("Tree", ref.Ref{})), false},
				},
				types.Choices{},
			),
		},
	}.New()
	return types.RegisterPackage(&p)
}

// Tree

type Tree struct {
	m types.Map
}

func NewTree() Tree {
	return Tree{types.NewMap(
		types.NewString("$name"), types.NewString("Tree"),
		types.NewString("$type"), types.MakeTypeRef("Tree", __testPackageInFile_struct_recursive_CachedRef),
		types.NewString("children"), types.NewList(),
	)}
}

type TreeDef struct {
	Children ListOfTreeDef
}

func (def TreeDef) New() Tree {
	return Tree{
		types.NewMap(
			types.NewString("$name"), types.NewString("Tree"),
			types.NewString("$type"), types.MakeTypeRef("Tree", __testPackageInFile_struct_recursive_CachedRef),
			types.NewString("children"), def.Children.New().NomsValue(),
		)}
}

func (s Tree) Def() (d TreeDef) {
	d.Children = ListOfTreeFromVal(s.m.Get(types.NewString("children"))).Def()
	return
}

var __typeRefForTree = types.MakeTypeRef("Tree", __testPackageInFile_struct_recursive_CachedRef)

func (m Tree) TypeRef() types.TypeRef {
	return __typeRefForTree
}

func TreeFromVal(val types.Value) Tree {
	// TODO: Validate here
	return Tree{val.(types.Map)}
}

func (s Tree) NomsValue() types.Value {
	return s.m
}

// A Noms Value that describes Tree.
var __typeRefForTree = types.MakeStructTypeRef("Tree",
	[]types.Field{
		types.Field{"children", types.MakeCompoundTypeRef("", types.ListKind, types.MakeTypeRef("Tree", ref.Ref{})), false},
	},
	types.Choices{},
)

func (m Tree) TypeRef() types.TypeRef {
	return __typeRefForTree
}

func init() {
	types.RegisterFromValFunction(__typeRefForTree, func(v types.Value) types.NomsValue {
		return TreeFromVal(v)
	})
}

func (s Tree) Equals(other Tree) bool {
	return s.m.Equals(other.m)
}

func (s Tree) Ref() ref.Ref {
	return s.m.Ref()
}

func (s Tree) Children() ListOfTree {
	return ListOfTreeFromVal(s.m.Get(types.NewString("children")))
}

func (s Tree) SetChildren(val ListOfTree) Tree {
	return Tree{s.m.Set(types.NewString("children"), val.NomsValue())}
}

// ListOfTree

type ListOfTree struct {
	l types.List
}

func NewListOfTree() ListOfTree {
	return ListOfTree{types.NewList()}
}

type ListOfTreeDef []TreeDef

func (def ListOfTreeDef) New() ListOfTree {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = d.New().NomsValue()
	}
	return ListOfTree{types.NewList(l...)}
}

func (l ListOfTree) Def() ListOfTreeDef {
	d := make([]TreeDef, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = TreeFromVal(l.l.Get(i)).Def()
	}
	return d
}

func ListOfTreeFromVal(val types.Value) ListOfTree {
	// TODO: Validate here
	return ListOfTree{val.(types.List)}
}

func (l ListOfTree) NomsValue() types.Value {
	return l.l
}

func (l ListOfTree) Equals(p ListOfTree) bool {
	return l.l.Equals(p.l)
}

func (l ListOfTree) Ref() ref.Ref {
	return l.l.Ref()
}

// A Noms Value that describes ListOfTree.
var __typeRefForListOfTree = types.MakeCompoundTypeRef("", types.ListKind, types.MakeTypeRef("Tree", ref.Ref{}))

func (m ListOfTree) TypeRef() types.TypeRef {
	return __typeRefForListOfTree
}

func init() {
	types.RegisterFromValFunction(__typeRefForListOfTree, func(v types.Value) types.NomsValue {
		return ListOfTreeFromVal(v)
	})
}

func (l ListOfTree) Len() uint64 {
	return l.l.Len()
}

func (l ListOfTree) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfTree) Get(i uint64) Tree {
	return TreeFromVal(l.l.Get(i))
}

func (l ListOfTree) Slice(idx uint64, end uint64) ListOfTree {
	return ListOfTree{l.l.Slice(idx, end)}
}

func (l ListOfTree) Set(i uint64, val Tree) ListOfTree {
	return ListOfTree{l.l.Set(i, val.NomsValue())}
}

func (l ListOfTree) Append(v ...Tree) ListOfTree {
	return ListOfTree{l.l.Append(l.fromElemSlice(v)...)}
}

func (l ListOfTree) Insert(idx uint64, v ...Tree) ListOfTree {
	return ListOfTree{l.l.Insert(idx, l.fromElemSlice(v)...)}
}

func (l ListOfTree) Remove(idx uint64, end uint64) ListOfTree {
	return ListOfTree{l.l.Remove(idx, end)}
}

func (l ListOfTree) RemoveAt(idx uint64) ListOfTree {
	return ListOfTree{(l.l.RemoveAt(idx))}
}

func (l ListOfTree) fromElemSlice(p []Tree) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v.NomsValue()
	}
	return r
}

type ListOfTreeIterCallback func(v Tree, i uint64) (stop bool)

func (l ListOfTree) Iter(cb ListOfTreeIterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(TreeFromVal(v), i)
	})
}

type ListOfTreeIterAllCallback func(v Tree, i uint64)

func (l ListOfTree) IterAll(cb ListOfTreeIterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(TreeFromVal(v), i)
	})
}

type ListOfTreeFilterCallback func(v Tree, i uint64) (keep bool)

func (l ListOfTree) Filter(cb ListOfTreeFilterCallback) ListOfTree {
	nl := NewListOfTree()
	l.IterAll(func(v Tree, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}
