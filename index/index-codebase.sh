#!/bin/bash

if [ "$1" == "" ]; then
  echo "$0 <directory>"
  exit -1
fi

# one container per package
for d in `find $1 -mindepth 2 -maxdepth 2 -type d`; do
  # $d is full path
  LNAME=`echo $d | sed -e 's/^\///' | sed -e 's/:\|\//_/g'`
  JSON_PATH=$1/$LNAME.json
  echo $JSON_PATH
  docker run --rm --name $LNAME --network="ctagsweb_default" -v "$PWD:/ctags-web" -v "$1:$1" -d ctags-web/web bash -c "ctags --recurse=yes --fields=* --output-format=json -f $JSON_PATH $d &&  /ctags-web/index/index db $JSON_PATH && rm -f $JSON_PATH"
  sleep $[ ( $RANDOM % 2 )  + 1 ]s
done
