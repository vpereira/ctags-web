#!/bin/bash
for d in import-go index-go web-go; do
  cd $d && go get -d -v ./... && go build && cd ..
done
