#!/bin/bash

GO_SRC=http://golang.org/dl/go1.3.src.tar.gz

set -e
set -x
#rm -fr ../tmp/dist/go
mkdir -p ../tmp/dist && cd ../tmp/dist

wget -O - ${GO_SRC} | tar xzf -

cd go/src
./make.bash
