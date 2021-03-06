// This file was generated by nomdl/codegen.
// @flow
/* eslint-disable */

import {
  Field as _Field,
  Kind as _Kind,
  Package as _Package,
  createStructClass as _createStructClass,
  float32Type as _float32Type,
  makeCompoundType as _makeCompoundType,
  makeListType as _makeListType,
  makeSetType as _makeSetType,
  makeStructType as _makeStructType,
  makeType as _makeType,
  newList as _newList,
  newSet as _newSet,
  registerPackage as _registerPackage,
  stringType as _stringType,
} from '@attic/noms';
import type {
  NomsList as _NomsList,
  NomsSet as _NomsSet,
  RefValue as _RefValue,
  Struct as _Struct,
  float32 as _float32,
} from '@attic/noms';

const _pkg = new _Package([
  _makeStructType('StructWithRef',
    [
      new _Field('r', _makeCompoundType(_Kind.Ref, _makeCompoundType(_Kind.Set, _float32Type)), false),
    ],
    [

    ]
  ),
], [
]);
_registerPackage(_pkg);
const StructWithRef$type = _makeType(_pkg.ref, 0);
const StructWithRef$typeDef = _makeStructType('StructWithRef',
  [
    new _Field('r', _makeCompoundType(_Kind.Ref, _makeCompoundType(_Kind.Set, _float32Type)), false),
  ],
  [

  ]
);


type StructWithRef$Data = {
  r: _RefValue<_NomsSet<_float32>>;
};

interface StructWithRef$Interface extends _Struct {
  constructor(data: StructWithRef$Data): void;
  r: _RefValue<_NomsSet<_float32>>;  // readonly
  setR(value: _RefValue<_NomsSet<_float32>>): StructWithRef$Interface;
}

export const StructWithRef: Class<StructWithRef$Interface> = _createStructClass(StructWithRef$type, StructWithRef$typeDef);

export function newListOfRefOfFloat32(values: Array<_RefValue<_float32>>): Promise<_NomsList<_RefValue<_float32>>> {
  return _newList(values, _makeListType(_makeCompoundType(_Kind.Ref, _float32Type)));
}

export function newListOfString(values: Array<string>): Promise<_NomsList<string>> {
  return _newList(values, _makeListType(_stringType));
}

export function newSetOfFloat32(values: Array<_float32>): Promise<_NomsSet<_float32>> {
  return _newSet(values, _makeSetType(_float32Type));
}
