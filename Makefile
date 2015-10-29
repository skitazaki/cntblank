all: clean build test dist

setup:
	go get -d

build: main.go app.go report.go dtparse.go
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
