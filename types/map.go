package types

type Map interface {
	Value
	First() (Value, Value)
	Len() uint64
	Empty() bool
	Has(key Value) bool
	Get(key Value) Value
	MaybeGet(key Value) (v Value, ok bool)
	Set(key Value, val Value) Map
	SetM(kv ...Value) Map
	Remove(k Value) Map
	Iter(cb mapIterCallback)
	IterAll(cb mapIterAllCallback)
	IterAllP(concurrency int, f mapIterAllCallback)
	Filter(cb mapFilterCallback) Map
	elemTypes() []Type
}

type indexOfMapFn func(m mapData, v Value) int
type mapIterCallback func(key, value Value) (stop bool)
type mapIterAllCallback func(key, value Value)
type mapFilterCallback func(key, value Value) (keep bool)

var mapType = MakeCompoundType(MapKind, MakePrimitiveType(ValueKind), MakePrimitiveType(ValueKind))

func NewMap(kv ...Value) Map {
	return NewTypedMap(mapType, kv...)
}

func NewTypedMap(t Type, kv ...Value) Map {
	return newTypedMap(t, buildMapData(mapData{}, kv, t)...)
}

func newTypedMap(t Type, entries ...mapEntry) Map {
	seq := newEmptySequenceChunker(makeMapLeafChunkFn(t, nil), newOrderedMetaSequenceChunkFn(t, nil), newMapLeafBoundaryChecker(), newOrderedMetaSequenceBoundaryChecker)

	for _, entry := range entries {
		seq.Append(entry)
	}

	return seq.Done().(Map)
}
