#!/bin/sh
set -e

rm -rf app
(cd server && gopherjs build -o ../app/server.js -m)
node_modules/.bin/webpack -p --progress
rm app/*.map
cp -v static/* app/
