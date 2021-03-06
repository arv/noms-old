package types

// NomsKind allows a TypeDesc to indicate what kind of type is described.
type NomsKind uint8

// All supported kinds of Noms types are enumerated here.
const (
	BoolKind NomsKind = iota
	Uint8Kind
	Uint16Kind
	Uint32Kind
	Uint64Kind
	Int8Kind
	Int16Kind
	Int32Kind
	Int64Kind
	Float32Kind
	Float64Kind
	StringKind
	BlobKind
	ValueKind
	ListKind
	MapKind
	RefKind
	SetKind
	EnumKind
	StructKind
	TypeKind
	UnresolvedKind
	PackageKind
)

// IsPrimitiveKind returns true if k represents a Noms primitive type, which excludes collections (List, Map, Set), Refs, Enums, Structs, Symbolic and Unresolved types.
func IsPrimitiveKind(k NomsKind) bool {
	switch k {
	case BoolKind, Int8Kind, Int16Kind, Int32Kind, Int64Kind, Float32Kind, Float64Kind, Uint8Kind, Uint16Kind, Uint32Kind, Uint64Kind, StringKind, BlobKind, ValueKind, TypeKind, PackageKind:
		return true
	default:
		return false
	}
}
