// @flow

export type MaybePromise<T> = T | Promise<T>;

type GeneratorObject = {
  next: (v: any) => ResultObject,
  throw: (ex: any) => ResultObject
};

type ResultObject = {
  done: boolean,
  value: any
};

export function spawn<T>(f: () => GeneratorObject): MaybePromise<T> {
  try {
    let gen = f();
    return step(gen, 'next', undefined);
  } catch (ex) {
    return Promise.reject(ex);
  }
}

function step<T>(gen: GeneratorObject, m: string, arg: any): MaybePromise<T> {
  let value, done;
  try {
    let res = gen[m](arg);
    value = res.value;
    done = res.done;
  } catch (ex) {
    return Promise.reject(ex);
  }
  if (done) {
    return value;
  }
  if (value instanceof Promise) {
    return value.then(
        v => step(gen, 'next', v),
        err => step(gen, 'throw', err));
  }
  return step(gen, 'next', value);
}
