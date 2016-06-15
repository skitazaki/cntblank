PROGRAM = cntblank

all: clean build test version dist

setup:
	go get github.com/constabulary/gb/...
	go get golang.org/x/tools/cmd/goimports
	go get github.com/jteeuwen/go-bindata/...
	gb vendor restore

build: src/cntblank/main.go src/cntblank/app.go src/cntblank/report.go src/cntblank/dtparse.go
	go-bindata -o src/${PROGRAM}/bindata.go templates
	go fmt src/${PROGRAM}/*
	go vet src/${PROGRAM}/*
	goimports -w src/${PROGRAM}/*
	gb build

test:
	gb test

version: build
	@./bin/cntblank --version
	@./bin/cntblank --help || :

local: build
	./bin/cntblank --output-format=html --output=_t.html testdata/addrcode_jp.xlsx

dist:
	env GOOS=darwin  GOARCH=386   gb build
	env GOOS=darwin  GOARCH=amd64 gb build
	env GOOS=linux   GOARCH=386   gb build
	env GOOS=linux   GOARCH=amd64 gb build
	env GOOS=windows GOARCH=386   gb build
	env GOOS=windows GOARCH=amd64 gb build
	-@md5sum bin/*

clean:
	rm -fr bin pkg

.PHONY: clean
