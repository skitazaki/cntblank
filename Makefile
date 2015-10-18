all: clean build test dist

setup:
	go get -d

build: app.go main.go
	go fmt
	go vet
	go build -o cntblank

test: app_test.go
	go test

dist:
	./build.sh

clean:
	rm -fr dist

.PHONY: clean
