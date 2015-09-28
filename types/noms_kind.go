package types

// NomsKind allows a TypeDesc to indicate what kind of type is described.
type NomsKind uint8

// All supported kinds of Noms types are enumerated here.
const (
	BoolKind NomsKind = iota
	UInt8Kind
	UInt16Kind
	UInt32Kind
	UInt64Kind
	Int8Kind
	Int16Kind
	Int32Kind   // 8
	Int64Kind   // 9
	Float32Kind // 10
	Float64Kind // 11
	StringKind  // 12
	BlobKind    // 13
	ValueKind   // 14
	ListKind    // 15
	MapKind     // 16
	RefKind     // 17
	SetKind     // 18
	EnumKind    // 19
	StructKind  // 20
	TypeRefKind
)

func IsPrimitiveKind(k NomsKind) bool {
	switch k {
	case BoolKind, Int8Kind, Int16Kind, Int32Kind, Int64Kind, Float32Kind, Float64Kind, UInt8Kind, UInt16Kind, UInt32Kind, UInt64Kind, StringKind, BlobKind, ValueKind, TypeRefKind:
		return true
	default:
		return false
	}
}
