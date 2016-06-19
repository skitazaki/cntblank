PROGRAM = cntblank

# Thanks - Crosscompiling Go applications with Make
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/
PLATFORMS := linux/386 linux/amd64 darwin/386 darwin/amd64 windows/386 windows/amd64
temp = $(subst /, , $@)
os   = $(word 1, $(temp))
arch = $(word 2, $(temp))

all: clean build test version local dist

setup:  ## Install development tools and libraries
	go get github.com/constabulary/gb/...
	go get golang.org/x/tools/cmd/goimports
	go get github.com/jteeuwen/go-bindata/...
	gb vendor restore

build: src/cntblank/main.go src/cntblank/app.go src/cntblank/report.go  ## Build binary after linting source files
	go-bindata -o src/${PROGRAM}/bindata.go templates
	go fmt src/${PROGRAM}/*
	go vet src/${PROGRAM}/*
	goimports -w src/${PROGRAM}/*
	gb build

test:  ## Run the unit tests
	gb test

version: build  ## Show version number and application usage
	@./bin/cntblank --version
	@./bin/cntblank --help || :

local: build  ## Run some tests on local machine
	@[ -d _build ] || mkdir _build
	./bin/cntblank --output-without-header --output-meta --output-format=csv --output=_build/t.txt testdata/addrcode_jp.xlsx testdata/prefecture_jp.tsv
	./bin/cntblank --output-format=json --output=_build/t.json testdata/addrcode_jp.xlsx
	./bin/cntblank --output-format=html --output=_build/t.html testdata/addrcode_jp.xlsx --input-delimiter=, testdata/elementary_school_teacher_ja.csv

dist: clean darwin linux windows  ## Build distribution binaries

darwin: darwin/386 darwin/amd64

linux: linux/386 linux/amd64

windows: windows/386 windows/amd64

$(PLATFORMS):
	env GOOS=$(os) GOARCH=$(arch) gb build

help:  ## Show this messages
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean:
	rm -fr bin pkg _build

.PHONY: clean
