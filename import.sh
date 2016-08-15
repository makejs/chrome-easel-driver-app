#!/bin/sh
set -e

FILE=$1
TDIR=$(mktemp -d)

echo "Using $TDIR for temp"

if ! [ -f "$FILE" ]
then
	echo "Usage: $0 ./EaselDriver-0.2.6.pkg"
	echo ""
	echo "Download the Mac version (.pkg file) and specify the path as the first argument to this script"
	exit 1
fi

xar -x -C "$TDIR" -f "$FILE"

mkdir -p iris-lib

(cd iris-lib && gunzip <"$TDIR/IrisLib-0.2.6.pkg/Payload" | cpio -i)

rm -rf "$TDIR"

