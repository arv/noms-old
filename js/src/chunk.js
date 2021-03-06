// @flow

import Ref from './ref.js';
import {encode, decode} from './utf8.js';

export default class Chunk {
  data: Uint8Array;
  _ref: ?Ref;

  constructor(data: Uint8Array = new Uint8Array(0), ref: ?Ref) {
    this.data = data;
    this._ref = ref;
  }

  get ref(): Ref {
    return this._ref || (this._ref = Ref.fromData(this.data));
  }

  isEmpty(): boolean {
    return this.data.length === 0;
  }

  toString(): string {
    return decode(this.data);
  }

  static emptyChunk: Chunk;

  static fromString(s: string, ref: ?Ref): Chunk {
    return new Chunk(encode(s), ref);
  }
}

export const emptyChunk = new Chunk();
