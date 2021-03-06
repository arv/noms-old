// This file was generated by nomdl/codegen.
// @flow
/* eslint-disable */

import {
  Field as _Field,
  Package as _Package,
  boolType as _boolType,
  createStructClass as _createStructClass,
  makeEnumType as _makeEnumType,
  makeStructType as _makeStructType,
  makeType as _makeType,
  registerPackage as _registerPackage,
  stringType as _stringType,
} from '@attic/noms';
import type {
  Struct as _Struct,
} from '@attic/noms';

const _pkg = new _Package([
  _makeStructType('S',
    [
      new _Field('s', _stringType, false),
      new _Field('b', _boolType, false),
    ],
    [

    ]
  ),
  _makeEnumType('E', 'e1', 'e2', 'e3'),
], [
]);
_registerPackage(_pkg);
const S$type = _makeType(_pkg.ref, 0);
const S$typeDef = _makeStructType('S',
  [
    new _Field('s', _stringType, false),
    new _Field('b', _boolType, false),
  ],
  [

  ]
);
const E$type = _makeType(_pkg.ref, 1);
const E$typeDef = _makeEnumType('E', 'e1', 'e2', 'e3');


type S$Data = {
  s: string;
  b: boolean;
};

interface S$Interface extends _Struct {
  constructor(data: S$Data): void;
  s: string;  // readonly
  setS(value: string): S$Interface;
  b: boolean;  // readonly
  setB(value: boolean): S$Interface;
}

export const S: Class<S$Interface> = _createStructClass(S$type, S$typeDef);

export type E =
  0 |  // e1
  1 |  // e2
  2;  // e3
