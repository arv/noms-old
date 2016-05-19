// @flow

import {invariant} from './assert.js';
import Struct from './struct.js';
import type {valueOrPrimitive} from './value.js';
import type RefValue from './ref-value.js';
import Set from './set.js';
import {Kind} from './noms-kind.js';
import type {CompoundDesc, StructDesc, Type} from './type.js';
import {
  getTypeOfValue,
  makeRefType,
  makeSetType,
  makeStructType,
  makeUnionType,
  valueType,
} from './type.js';
import {equals} from './compare.js';

export default class Commit<T: valueOrPrimitive> extends Struct {
  constructor(value: T, parents: Set<RefValue<Commit>> = new Set()) {
    const commitType = findCommitType(parents.type, getTypeOfValue(value));
    super(commitType, {value, parents});
  }

  get value(): T {
    // $FlowIssue: _data is private.
    const value: T = this._data.value;
    return value;
  }

  setValue<U: valueOrPrimitive>(value: U): Commit<U> {
    return new Commit(value, this.parents);
  }

  get parents(): Set<RefValue<Commit>> {
    // $FlowIssue: _data is private.
    const parents: Set<RefValue<Commit>> = this._data.parents;
    invariant(parents instanceof Set);
    return parents;
  }

  setParents(parents: Set<RefValue<Commit>>): Commit<T> {
    return new Commit(this.value, parents);
  }
}
//
// export function newCommit(value: valueOrPrimitive,
//                           parentsArr: Array<RefValue<Commit>> = []): Promise<Commit {
//   const st = findCommitType(parentsArr.map(v => v.type), getTypeOfValue(value));
//   return new Set(parentsArr).then(parents => newStructWithType(st, {value, parents}));
// }

function unpackCompoundType(t: Type<CompoundDesc>): Type {
  return t.desc.elemTypes[0];
}

// Ref<T> -> T
const unpackRefType = unpackCompoundType;
// Set<T> -> T
const unpackSetType = unpackCompoundType;

function unpackCommitValueType(type: Type<StructDesc>): Type {
  return type.desc.fields['value'];
}

function unpackMaybeUnionType(t: Type): Type[] {
  if (t.kind === Kind.Union) {
    return t.desc.elemTypes;
  }
  return [t];
}

// export for testing
export function findCommitType(parentSetType: Type<CompoundDesc>, vt: Type): Type<StructDesc> {
  const parentTypes = unpackMaybeUnionType(unpackSetType(parentSetType));
  for (let i = 0; i < parentTypes.length; i++) {
    const pst = unpackRefType(parentTypes[i]);
    const pvt = unpackCommitValueType(pst);
    if (equals(vt, pvt)) {
      return pst;
    }
  }

  const st = makeStructType('Commit', {
    value: vt,
    parents: valueType,  // placeholder
  });
  parentTypes.push(makeRefType(st));
  st.desc.fields['parents'] = makeSetType(makeUnionType(parentTypes));
  return st;
}
