package types

import "crypto/sha1"

func newIndexedMetaSequenceBoundaryChecker() boundaryChecker {
	return newBuzHashBoundaryChecker(objectWindowSize, sha1.Size, objectPattern, func(item sequenceItem) []byte {
		digest := item.(metaTuple).ChildRef().TargetRef().Digest()
		return digest[:]
	})
}

// If |sink| is not nil, chunks will be eagerly written as they're created. Otherwise they are
// written when the root is written.
func newIndexedMetaSequenceChunkFn(t Type, source ValueReader, sink ValueWriter) makeChunkFn {
	return func(items []sequenceItem) (sequenceItem, Value) {
		tuples := make(metaSequenceData, len(items))

		for i, v := range items {
			tuples[i] = v.(metaTuple)
		}

		meta := newMetaSequenceFromData(tuples, t, source)
		if sink != nil {
			r := sink.WriteValue(meta)
			return newMetaTuple(Uint64(tuples.uint64ValuesSum()), nil, r), meta
		}
		return newMetaTuple(Uint64(tuples.uint64ValuesSum()), meta, Ref{}), meta
	}
}
