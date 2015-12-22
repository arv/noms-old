// @flow

type MaybePromise<T> = T | Promise<T>;

export function spawn<T>(f: () => any): MaybePromise<T> {
  try {
    let gen = f();
    return step(gen, 'next', undefined);
  } catch (ex) {
    return Promise.reject(ex);
  }
}

function step(gen, m: string, arg: any) {
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
