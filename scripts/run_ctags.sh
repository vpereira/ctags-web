# relative to the root directory
DIR="testProject"
ctags --recurse=yes --fields=* --output-format=json -f ctags.json $DIR
