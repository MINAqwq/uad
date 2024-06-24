#!/bin/sh

GO=go
BIN=uad
BUILD_ARGS=-v

if [ -z $1 ]; then
    echo "$0 [source_dir]"
    exit 1
fi

cd $1;
echo "[BUILD] $BIN"
$GO build $BUILD_ARGS .
cd -

if [ -z $1/$BIN ]; then
	echo "[BUILD] build failed!"
	exit 1
fi

mkdir -v -p build
mv -v $1/$BIN build

