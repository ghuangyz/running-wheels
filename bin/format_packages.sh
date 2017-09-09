#!/bin/bash

for pkg in `find ./pkg/* -type d`; do
    pushd $pkg
    go fmt
    popd
done
