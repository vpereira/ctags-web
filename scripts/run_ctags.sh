#!/bin/bash
# relative to the root directory
if [ "$1" == "" ]; then
  DIR="testProject"
else
  DIR=$1
fi

ctags --recurse=yes --fields=* --output-format=json -f ctags.json $DIR
