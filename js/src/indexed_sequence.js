// @flow

import type {ChunkStore} from './chunk_store.js';
import {notNull} from './assert.js';
import {search, Sequence, SequenceCursor} from './sequence.js';
import {spawn} from './spawn.js';
import type {MaybePromise} from './spawn.js';

export class IndexedSequence<T> extends Sequence<T> {
  getOffset(idx: number): number { // eslint-disable-line no-unused-vars
    throw new Error('override');
  }

  newCursorAt(cs: ChunkStore, idx: number): MaybePromise<IndexedSequenceCursor> {
    let cursor: ?IndexedSequenceCursor = null;
    let sequence: ?IndexedSequence = this;

    return spawn(function*() {
      while (sequence) {
        cursor = new IndexedSequenceCursor(cs, cursor, sequence, 0);
        idx -= cursor.advanceToOffset(idx);
        sequence = yield cursor.getChildSequence();
      }

      return notNull(cursor);
    });
  }
}

export class IndexedSequenceCursor<T> extends SequenceCursor<T, IndexedSequence> {
  advanceToOffset(idx: number): number {
    this.idx = search(this.length, (i: number) => {
      return idx <= this.sequence.getOffset(i);
    });

    if (this.idx === this.length) {
      this.idx = this.length - 1;
    }

    return this.idx > 0 ? this.sequence.getOffset(this.idx - 1) + 1 : 0;
  }
}
