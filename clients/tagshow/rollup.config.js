const assert = require('assert');
const babel = require('rollup-plugin-babel');
const commonjs = require('rollup-plugin-commonjs');
const inject = require('rollup-plugin-inject');
const npm = require('rollup-plugin-npm');
const path = require('path');
const replace = require('rollup-plugin-replace');

module.exports = {
  plugins: [
    babel(),
    inject({
      regeneratorRuntime: 'babel-regenerator-runtime'
    }),
    replaceEnv('NOMS_SERVER'),
    {
      resolveId(importee, importer) {
        // TODO: Find browser entry in package.json?
        if (importee === './fetch.js') {
          return path.resolve(path.dirname(importer), 'browser', 'fetch.js');
        }
      }
    },
    npm({
      jsnext: true,
      main: true,
      skip: []
    }),
    commonjs()
  ],
  format: 'iife'
};

function replaceEnv(/*...names*/) {
  let o = {};
  let names = ['NODE_ENV', ...arguments];
  for (let n of names) {
    assert(n in process.env, `Missing environment variable: ${n}`);
    o[`process.env.${n}`] = JSON.stringify(process.env[n]);
  }
  return replace(o);
}
