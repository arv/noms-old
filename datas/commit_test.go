package datas

import (
	"testing"

	"github.com/attic-labs/noms/types"
	"github.com/stretchr/testify/assert"
)

func TestCommitType(t *testing.T) {
	assert := assert.New(t)

	a1 := findCommitType(nil, types.NumberType)
	t1 := types.MakeStructType("Commit", types.TypeMap{
		"value":   types.NumberType,
		"parents": types.ValueType, // placeholder
	})
	t1.Desc.(types.StructDesc).Fields["parents"] = types.MakeSetType(types.MakeRefType(t1))
	assert.True(a1.Equals(t1))

	a2 := findCommitType([]*types.Type{types.MakeRefType(t1)}, types.NumberType)
	assert.True(a2.Equals(t1))

	a3 := findCommitType([]*types.Type{types.MakeRefType(t1)}, types.StringType)
	t2 := types.MakeStructType("Commit", types.TypeMap{
		"value":   types.StringType,
		"parents": types.ValueType, // placeholder
	})
	t2.Desc.(types.StructDesc).Fields["parents"] = types.MakeSetType(types.MakeUnionType(types.MakeRefType(t1), types.MakeRefType(t2)))
	assert.True(a3.Equals(t2))

	a4 := findCommitType([]*types.Type{types.MakeRefType(t1), types.MakeRefType(t2)}, types.StringType)
	assert.True(a4.Equals(t2))

	a5 := findCommitType([]*types.Type{types.MakeRefType(t1), types.MakeRefType(t2)}, types.NumberType)
	assert.True(a5.Equals(t1))

}
