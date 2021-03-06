package types

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/attic-labs/noms/ref"
	"github.com/stretchr/testify/assert"
)

func TestWriteValue(t *testing.T) {
	assert := assert.New(t)

	testEncode := func(expected string, v Value) ref.Ref {
		vs := NewTestValueStore()
		r := vs.WriteValue(v).TargetRef()

		// Assuming that MemoryStore works correctly, we don't need to check the actual serialization, only the hash. Neat.
		assert.EqualValues(sha1.Sum([]byte(expected)), r.Digest(), "Incorrect ref serializing %+v. Got: %#x", v, r.Digest())
		return r
	}

	// Encoding details for each codec is tested elsewhere.
	// Here we just want to make sure codecs are selected correctly.
	b := NewBlob(bytes.NewBuffer([]byte{0x00, 0x01, 0x02}))
	testEncode(string([]byte{'b', ' ', 0x00, 0x01, 0x02}), b)

	testEncode(fmt.Sprintf("t [%d,\"hi\"]", StringKind), NewString("hi"))

	testEncode(fmt.Sprintf("t [%d,[],[]]", PackageKind), Package{types: []Type{}, dependencies: []ref.Ref{}, ref: &ref.Ref{}})
	ref1 := testEncode(fmt.Sprintf("t [%d,[%d],[]]", PackageKind, BoolKind), Package{types: []Type{MakePrimitiveType(BoolKind)}, dependencies: []ref.Ref{}, ref: &ref.Ref{}})
	testEncode(fmt.Sprintf("t [%d,[],[\"%s\"]]", PackageKind, ref1), Package{types: []Type{}, dependencies: []ref.Ref{ref1}, ref: &ref.Ref{}})
}

func TestWriteBlobLeaf(t *testing.T) {
	assert := assert.New(t)
	vs := NewTestValueStore()

	buf := bytes.NewBuffer([]byte{})
	b1 := NewBlob(buf)
	bl1, ok := b1.(blobLeaf)
	assert.True(ok)
	r1 := vs.WriteValue(bl1).TargetRef()
	// echo -n 'b ' | sha1sum
	assert.Equal("sha1-e1bc846440ec2fb557a5a271e785cd4c648883fa", r1.String())

	buf = bytes.NewBufferString("Hello, World!")
	b2 := NewBlob(buf)
	bl2, ok := b2.(blobLeaf)
	assert.True(ok)
	r2 := vs.WriteValue(bl2).TargetRef()
	// echo -n 'b Hello, World!' | sha1sum
	assert.Equal("sha1-135fe1453330547994b2ce8a1b238adfbd7df87e", r2.String())
}

func TestWritePackageWhenValueIsWritten(t *testing.T) {
	assert := assert.New(t)
	vs := NewTestValueStore()

	typeDef := MakeStructType("S", []Field{
		Field{"X", MakePrimitiveType(Int32Kind), false},
	}, Choices{})
	pkg1 := NewPackage([]Type{typeDef}, []ref.Ref{})
	// Don't write package
	pkgRef1 := RegisterPackage(&pkg1)
	typ := MakeType(pkgRef1, 0)

	s := NewStruct(typ, typeDef, structData{"X": Int32(42)})
	vs.WriteValue(s)

	pkg2 := vs.ReadValue(pkgRef1)
	assert.True(pkg1.Equals(pkg2))
}

func TestWritePackageDepWhenPackageIsWritten(t *testing.T) {
	assert := assert.New(t)
	vs := NewTestValueStore()

	pkg1 := NewPackage([]Type{}, []ref.Ref{})
	// Don't write package
	pkgRef1 := RegisterPackage(&pkg1)

	pkg2 := NewPackage([]Type{}, []ref.Ref{pkgRef1})
	vs.WriteValue(pkg2)

	pkg3 := vs.ReadValue(pkgRef1)
	assert.True(pkg1.Equals(pkg3))
}
