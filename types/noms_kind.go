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
	Int32Kind
	Int64Kind   // 8
	Float32Kind // 9
	Float64Kind // 10
	StringKind  // 11
	BlobKind    // 12
	ValueKind   // 13
	ListKind    // 14
	MapKind     // 15
	RefKind     // 16
	SetKind     // 17
	EnumKind    // 18
	StructKind  // 19
	TypeRefKind // 20
)

func IsPrimitiveKind(k NomsKind) bool {
	switch k {
	case BoolKind, Int8Kind, Int16Kind, Int32Kind, Int64Kind, Float32Kind, Float64Kind, UInt8Kind, UInt16Kind, UInt32Kind, UInt64Kind, StringKind, BlobKind, ValueKind, TypeRefKind:
		return true
	default:
		return false
	}
}
