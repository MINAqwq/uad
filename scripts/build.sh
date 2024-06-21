#!/bin/sh

if [ -z $1 ]; then
    echo "$0 [source_dir]"
    exit 1
fi

cd $1;
go build -v .
cd -;
