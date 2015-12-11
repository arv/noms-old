#!/bin/bash
SRC="main.js"
OUT="out.js"

NOMS_ROOT="../.."
NOMS_SRC="$NOMS_ROOT/js/src"
ROLLUP_CMD="node_modules/.bin/rollup --config --input $SRC"

export NODE_ENV=$1
export BABEL_ENV=$1

if [ $1 == "production" ]; then
  $ROLLUP_CMD | node_modules/.bin/uglifyjs -c warnings=false -m -w --screw-ie8 > $OUT
elif [ $1 == "development" ]; then
  ROLLUP_CMD="$ROLLUP_CMD --sourcemap inline"

  node_modules/.bin/http-server &
  node_modules/.bin/nodemon --watch . --watch $NOMS_SRC --ignore $OUT \
      -x "$ROLLUP_CMD > $OUT"
fi
