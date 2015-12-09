SRC="main.js"
OUT="out.js"

export NODE_ENV=production
export BABEL_ENV=production-rollup

node_modules/.bin/rollup --config --input $SRC \
    | node_modules/.bin/uglifyjs -c warnings=false -m -w --screw-ie8 > $OUT
