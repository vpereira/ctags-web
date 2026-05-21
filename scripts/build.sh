#!/bin/bash
# Build all ctags-web binaries
set -e

cd "$(dirname "$0")/.."

echo "==> Downloading dependencies..."
go mod tidy

echo "==> Building index..."
go build -o index/ctags-index ./index/

echo "==> Building import..."
go build -o import/ctags-import ./import/

echo "==> Building web..."
go build -o web/ctags-web ./web/

echo "==> Done"
