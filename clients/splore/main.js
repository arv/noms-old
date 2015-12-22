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
  let r = Ref.parse('sha1-5380f69030ad209e137dee221704c34b2717a02b');  // Category Software
  let c = Ref.parse('sha1-c6b9977d49a1ee414c78f305fd0f923352864464');  // Time: Last year

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
