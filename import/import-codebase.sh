#!/bin/bash

if [ "$1" == "" ]; then
  echo "$0 <directory>"
  exit -1
fi

for d in `find $1 -maxdepth 1 -type d`; do
  docker run -ti --rm --network="ctagsweb_default" -v "$PWD:/ctags-web" -v "$1:$1" ctags-web/web bash -c "/ctags-web/import/import db 'ctags' 'code' $d && echo $?"
done
