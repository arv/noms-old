package csv

import (
	"math"
	"strconv"

	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/types"
)

type schemaOptions []*typeCanFit

func newSchemaOptions(fieldCount int) schemaOptions {
	options := make([]*typeCanFit, fieldCount, fieldCount)
	for i := 0; i < fieldCount; i++ {
		options[i] = &typeCanFit{true, true, true, true, true, true, true, true, true, true, true, true, true, true}
	}
	return options
}

func (so schemaOptions) Test(fields []string) {
	for i, t := range so {
		if i < len(fields) {
			t.Test(fields[i])
		}
	}
}

func (so schemaOptions) ValidKinds() []KindSlice {
	kinds := make([]KindSlice, len(so))
	for i, t := range so {
		kinds[i] = t.ValidKinds()
	}
	return kinds
}

type typeCanFit struct {
	uintType    bool
	intType     bool
	boolType    bool
	uint8Type   bool
	uint16Type  bool
	uint32Type  bool
	uint64Type  bool
	int8Type    bool
	int16Type   bool
	int32Type   bool
	int64Type   bool
	float32Type bool
	float64Type bool
	stringType  bool
}

func (tc *typeCanFit) ValidKinds() (kinds KindSlice) {
	if tc.uintType {
		if tc.uint8Type {
			kinds = append(kinds, types.Uint8Kind)
		}
		if tc.uint16Type {
			kinds = append(kinds, types.Uint16Kind)
		}
		if tc.uint32Type {
			kinds = append(kinds, types.Uint32Kind)
		}
		if tc.uint64Type {
			kinds = append(kinds, types.Uint64Kind)
		}
	}
	if tc.intType {
		if tc.int8Type {
			kinds = append(kinds, types.Int8Kind)
		}
		if tc.int16Type {
			kinds = append(kinds, types.Int16Kind)
		}
		if tc.int32Type {
			kinds = append(kinds, types.Int32Kind)
		}
		if tc.int64Type {
			kinds = append(kinds, types.Int64Kind)
		}
	}
	if tc.float32Type {
		kinds = append(kinds, types.Float32Kind)
	}
	if tc.float64Type {
		kinds = append(kinds, types.Float64Kind)
	}
	if tc.boolType {
		kinds = append(kinds, types.BoolKind)
	}

	kinds = append(kinds, types.StringKind)
	return kinds
}

func (tc *typeCanFit) Test(value string) {
	tc.testUints(value)
	tc.testInts(value)
	tc.testFloats(value)
	tc.testBool(value)
}

func (tc *typeCanFit) testUints(value string) {
	if !tc.uintType {
		return
	}

	ival, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		tc.uintType = false
		tc.uint8Type = false
		tc.uint16Type = false
		tc.uint32Type = false
		tc.uint64Type = false
		return
	}

	tc.uint32Type = tc.uint32Type && ival <= math.MaxUint32
	tc.uint16Type = tc.uint16Type && ival <= math.MaxUint16
	tc.uint8Type = tc.uint8Type && ival <= math.MaxUint8
	return
}

func (tc *typeCanFit) testInts(value string) {
	if !tc.intType {
		return
	}

	ival, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		tc.intType = false
		tc.int8Type = false
		tc.int16Type = false
		tc.int32Type = false
		tc.int64Type = false
		return
	}

	if ival < 0 {
		ival *= -1
	}
	tc.int32Type = tc.int32Type && ival <= math.MaxInt32
	tc.int16Type = tc.int16Type && ival <= math.MaxInt16
	tc.int8Type = tc.int8Type && ival <= math.MaxInt8
	return
}

func (tc *typeCanFit) testFloats(value string) {
	if !tc.float32Type && !tc.float64Type {
		return
	}

	fval, err := strconv.ParseFloat(value, 64)
	if err != nil {
		tc.float32Type = false
		tc.float64Type = false
		return
	}

	if fval > math.MaxFloat32 {
		tc.float32Type = false
	}
}

func (tc *typeCanFit) testBool(value string) {
	if !tc.boolType {
		return
	}
	_, err := strconv.ParseBool(value)
	tc.boolType = err == nil
}

// StringToType takes a piece of data as a string and attempts to convert it to a types.Value of the appropriate types.NomsKind.
func StringToType(s string, k types.NomsKind) types.Value {
	switch k {
	case types.Uint8Kind:
		ival, err := strconv.ParseUint(s, 10, 8)
		d.Chk.NoError(err)
		return types.Uint8(ival)
	case types.Uint16Kind:
		ival, err := strconv.ParseUint(s, 10, 16)
		d.Chk.NoError(err)
		return types.Uint16(ival)
	case types.Uint32Kind:
		ival, err := strconv.ParseUint(s, 10, 32)
		d.Chk.NoError(err)
		return types.Uint32(ival)
	case types.Uint64Kind:
		ival, err := strconv.ParseUint(s, 10, 64)
		d.Chk.NoError(err)
		return types.Uint64(ival)
	case types.Int8Kind:
		ival, err := strconv.ParseInt(s, 10, 8)
		d.Chk.NoError(err)
		return types.Int8(ival)
	case types.Int16Kind:
		ival, err := strconv.ParseInt(s, 10, 16)
		d.Chk.NoError(err)
		return types.Int16(ival)
	case types.Int32Kind:
		ival, err := strconv.ParseInt(s, 10, 32)
		d.Chk.NoError(err)
		return types.Int32(ival)
	case types.Int64Kind:
		ival, err := strconv.ParseInt(s, 10, 64)
		d.Chk.NoError(err)
		return types.Int64(ival)
	case types.Float32Kind:
		fval, err := strconv.ParseFloat(s, 32)
		d.Chk.NoError(err)
		return types.Float32(fval)
	case types.Float64Kind:
		fval, err := strconv.ParseFloat(s, 64)
		d.Chk.NoError(err)
		return types.Float64(fval)
	case types.BoolKind:
		bval, err := strconv.ParseBool(s)
		d.Chk.NoError(err)
		return types.Bool(bval)
	case types.StringKind:
		return types.NewString(s)
	default:
		d.Exp.Fail("Invalid column type kind:", k)
	}
	panic("not reached")
}
