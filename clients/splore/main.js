// @flow

import {clearPackages, HttpStore, readValue, NomsSet, Ref} from 'noms';

let httpStore: HttpStore;

const nomsServer: ?string = process.env.NOMS_SERVER;
if (!nomsServer) {
  throw new Error('NOMS_SERVER not set');
}

function loadAllRefs(r: Ref): Promise<void> {
  return new Promise(resolve => {
    let wg = 0;

    let step = (r: Ref) => {
      readValue(r, httpStore).then(v => {
        wg--;
        if (v instanceof NomsSet) {
          let seq = v.sequence;
          if (seq.isMeta) {
            for (let mt of seq.items) {
              wg++;
              step(mt.ref);
            }
          } else if (wg === 0) {
            resolve();
          }
        } else {
          throw new Error('ouch');
        }
      });
    };

    wg++;
    step(r);
  });
}

window.addEventListener('load', async () => {
  httpStore = new HttpStore(nomsServer);
  let r = Ref.parse('sha1-ed02320e8904294233627d4d44ec201148353c38');
  let c = Ref.parse('sha1-6f8b0db017c5d2c6c0436db5daea08012dbcc51c');

  await loadAllRefs(r);
  await loadAllRefs(c);

  let s1:NomsSet<Ref> = await readValue(r, httpStore);
  let s2:NomsSet<Ref> = await readValue(c, httpStore);
//   clearPackages();

  let t1 = performance.now();

  await s1.intersect(s2);

  let t2 = performance.now();

  document.getElementById('splore').innerHTML = (t2 - t1) + 'ms';
});
