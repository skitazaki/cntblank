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
	@[ -d _build ] || mkdir _build
	./bin/cntblank --output-without-header --output-meta --output-format=csv --output=_build/t.txt testdata/addrcode_jp.xlsx testdata/prefecture_jp.tsv
	./bin/cntblank --output-format=json --output=_build/t.json testdata/addrcode_jp.xlsx
	./bin/cntblank --output-format=html --output=_build/t.html testdata/addrcode_jp.xlsx --input-delimiter=, testdata/elementary_school_teacher_ja.csv

dist:
	env GOOS=darwin  GOARCH=386   gb build
	env GOOS=darwin  GOARCH=amd64 gb build
	env GOOS=linux   GOARCH=386   gb build
	env GOOS=linux   GOARCH=amd64 gb build
	env GOOS=windows GOARCH=386   gb build
	env GOOS=windows GOARCH=amd64 gb build
	-@md5sum bin/*

clean:
	rm -fr bin pkg _build

.PHONY: clean
