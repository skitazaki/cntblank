all: clean build test version dist

setup:
	go get github.com/constabulary/gb/...
	gb vendor restore

build: src/cntblank/main.go src/cntblank/app.go src/cntblank/report.go src/cntblank/dtparse.go
	go fmt src/cntblank/*
	go vet src/cntblank/*
	gb build

test:
	gb test

version: build
	@./bin/cntblank --version
	@./bin/cntblank --help || :

dist:
	env GOOS=darwin  GOARCH=386   gb build
	env GOOS=darwin  GOARCH=amd64 gb build
	env GOOS=linux   GOARCH=386   gb build
	env GOOS=linux   GOARCH=amd64 gb build
	env GOOS=windows GOARCH=386   gb build
	env GOOS=windows GOARCH=amd64 gb build
	md5sum bin/*

clean:
	rm -fr bin pkg

.PHONY: clean
