// @flow

'use strict';

function go(generator) {
  let gen = generator();
  step(gen, gen.next());
}

const CONTINUE = 0;
const EXCEPTION = 1;
const PROMISE = 2;

// type Task =
//     {kind: 0, value: Promise<any>} |
//     {kind: 1, value: any};

// function stepVal(gen, v) {
//   step(gen, {done: false, value: {kind: CONTINUE, value: v}});
// }
//
// function stepEx(gen, v) {
//   step(gen, {done: false, value: {kind: EXCEPTION, value: ex}});
// }

class Task {
  constructor(gen) {
    this.gen = gen;
  }
}

function step(gen, res) {
  while (!res.done) {
    let task = res.value;
    let val = task.value;
    switch (pair.kind) {
      case CONTINUE:
        res = getNext(gen, val);
        break;
      case EXCEPTION:
        res = getThrow(gen, val);
        break;
      case PROMISE:
        val.then(v => { step(gen, {done: false, value: {kind: CONTINUE, value: v}}); },
                 ex => { step(gen, {done: false, value: {kind: EXCEPTION, value: ex}}); });
        return;
    }
  }
}

function getNext(gen, val) {
  try {
    return gen.next(val);
  } catch (ex) {
    return {done: false, value: {kind: EXCEPTION, value: ex}};
  }
}

function getThrow(gen, val) {
  try {
    return gen.throw(val);
  } catch (ex) {
    return {done: false, value: {kind: EXCEPTION, value: ex}};
  }
}

function wait(ms) {
  return {kind: PROMISE, value: new Promise((res) => {
    setTimeout(() => {
      res('abc');
    }, ms);
  })};
}

const N = 1e6;

go(function*() {
  // let d = Date.now();
  for (let i = 0; i < N; i++) {
    let x = yield {kind: CONTINUE, value: i};
    /*process.stdout.write*/(`\r${x}`);
  }
  // console.log('\ngo', Date.now() - d);
});

// (async () => {
//   let d = Date.now();
//   for (let i = 0; i < N; i++) {
//     let x = await i;
//     /*process.stdout.write*/(`\r${x}`);
//   }
//   console.log('\nasync', Date.now() - d);
// })();
