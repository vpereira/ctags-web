#!/bin/bash

if [ "$1" == "" ]; then
  echo "$0 <directory>"
  exit -1
fi

# one container per codestream
for d in `find $1 -maxdepth 1 -type d`; do
  docker run -ti --rm --network="ctagsweb_default" -v "$PWD:/ctags-web" -v "$1:$1" ctags-web/web bash -c "ctags --recurse=yes --fields=* --output-format=json -f $d.json $d && \
                                /ctags-web/index/index db $d.json  && echo $?"
done
