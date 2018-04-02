#!/bin/bash

mongo --host db --eval 'rs.initiate()'
# waiting until es bootstrap
sleep 10
# maybe move it to its own container
mongo-connector -m db:27017 -t es:9200 -d elastic_doc_manager
# just to make it run forever
# it should start the server instead tail -f
tail -f /dev/null
