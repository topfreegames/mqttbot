#!/bin/bash

echo "mode: count" > profile.cov
for pkg in `cat testpackages.txt`
do
    go test -v -covermode=count -coverprofile=coverage.out $pkg
    tail -n +2 coverage.out >> profile.cov
done
