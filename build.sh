#!/bin/bash

for os in linux
do
    for arch in amd64
    do
        echo "go build -v -o bin/pa-go-starter-${arch}-${os}"
        GOOS=${os} GOARCH=${arch} go build -v -o bin/pa-go-starter-${arch}-${os}
    done
done