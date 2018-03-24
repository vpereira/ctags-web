all: deps
	go build

deps:
	go get -d -v ./...

tests: deps
	go test

.PHONY: clean
clean:
	rm web
