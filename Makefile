SUBDIRS := import index web

.PHONY: all clean build

all: build

build:
	@go mod tidy
	@go build -o index/ctags-index ./index/
	@go build -o import/ctags-import ./import/
	@go build -o web/ctags-web ./web/

.PHONY: clean
clean:
	rm -f index/ctags-index import/ctags-import web/ctags-web

.PHONY: test
test:
	go test -short ./...

.PHONY: vet
vet:
	go vet ./...
