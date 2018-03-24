all: deps
	go build

deps:
	go get -d -v ./...

tests: deps
	go test

clean:
	rm import
