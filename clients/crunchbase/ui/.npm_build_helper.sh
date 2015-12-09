#!/bin/bash
SRC="src/main.js"
OUT="out.js"

export NODE_ENV=production
export BABEL_ENV=production-rollup

cp node_modules/nvd3/build/nv.d3.min.css nvd3.css
cp node_modules/nvd3/build/nv.d3.min.js nvd3.js
cp node_modules/d3/d3.min.js d3.js

node_modules/.bin/rollup --config --input $SRC \
    | node_modules/.bin/uglifyjs -c warnings=false -m -w --screw-ie8 > $OUT
