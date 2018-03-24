#!/bin/bash
for d in import index web; do
  cd $d && go get -d -v ./... && go build && cd ..
done
