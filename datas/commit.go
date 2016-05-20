package datas

import (
	"fmt"

	"github.com/attic-labs/noms/types"
)

var commitType *types.Type

const (
	ParentsField = "parents"
	ValueField   = "value"
)

func init() {
	structName := "Commit"

	// struct Commit {
	//   parents: Set<Ref<Commit>>
	//   value: Value
	// }

	fieldTypes := types.TypeMap{
		ParentsField: nil,
		ValueField:   types.ValueType,
	}
	commitType = types.MakeStructType(structName, fieldTypes)
	commitType.Desc.(types.StructDesc).Fields[ParentsField] = types.MakeSetType(types.MakeRefType(commitType))
}

func NewCommit(value types.Value, parents ...types.Ref) types.Struct {
	parentValues := make([]types.Value, len(parents))
	parentTypes := make([]*types.Type, len(parents))
	for i, p := range parents {
		parentValues[i] = p
		parentTypes[i] = p.Type()
	}
	st := findCommitType(parentTypes, value.Type())
	initialFields := map[string]types.Value{
		ValueField:   value,
		ParentsField: types.NewSet(parentValues...),
	}
	return types.NewStructWithTypeNoValidation(st, initialFields)
}

func findCommitType(parentTypes []*types.Type, vt *types.Type) *types.Type {
	for _, pt := range parentTypes {
		pst := pt.Desc.(types.CompoundDesc).ElemTypes[0]
		pvt := pst.Desc.(types.StructDesc).Fields["value"]
		if vt.Equals(pvt) {
			return pst
		}
	}

	st := types.MakeStructType("Commit", types.TypeMap{
		"value":   vt,
		"parents": types.ValueType, // placeholder
	})
	parentTypes = append(parentTypes, types.MakeRefType(st))
	st.Desc.(types.StructDesc).Fields["parents"] = types.MakeSetType(types.MakeUnionType(parentTypes...))
	return st
}

func IsCommitType(t *types.Type) bool {
	fmt.Printf("\n\ncommitType:\n%s\n\nt:\n%s\n\n", commitType.Describe(), t.Describe())
	return types.IsSubtype(commitType, t)
}
