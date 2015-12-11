#!/bin/bash
SRC="src/main.js"
OUT="out.js"

NOMS_ROOT="../../.."
NOMS_SRC="$NOMS_ROOT/js/src"
ROLLUP_CMD="node_modules/.bin/rollup --config --input $SRC"

export NODE_ENV=$1
export BABEL_ENV=$1

if [ $1 == "production" ]; then
  cp node_modules/nvd3/build/nv.d3.min.css nvd3.css
  cp node_modules/nvd3/build/nv.d3.min.js nvd3.js
  cp node_modules/d3/d3.min.js d3.js

  $ROLLUP_CMD | node_modules/.bin/uglifyjs -c warnings=false -m -w --screw-ie8 > $OUT
elif [ $1 == "development" ]; then
  cp node_modules/nvd3/build/nv.d3.css nvd3.css
  cp node_modules/nvd3/build/nv.d3.js nvd3.js
  cp node_modules/d3/d3.js d3.js

  ROLLUP_CMD="$ROLLUP_CMD --sourcemap inline"

  node_modules/.bin/http-server &
  node_modules/.bin/nodemon --watch . --watch $NOMS_SRC --ignore $OUT \
      -x "$ROLLUP_CMD > $OUT"
fi
