#!/bin/sh

set -eu

APP=cntblank
DOCKER_IMAGE="golang:1.4-cross"
DOCKER_WORKDIR="/usr/src/$APP"

cat <<EOF >make.bash
#/bin/bash
set -eux

go get -d -v
# go test -v 2>&1 | tee test.txt
for GOOS in darwin linux; do
    for GOARCH in 386 amd64; do
        env GOOS=\$GOOS GOARCH=\$GOARCH go build -v -o $APP-\$GOOS-\$GOARCH
    done
done
for GOOS in windows; do
    for GOARCH in 386 amd64; do
        env GOOS=\$GOOS GOARCH=\$GOARCH go build -v -o $APP-\$GOOS-\$GOARCH.exe
    done
done
EOF
chmod +x make.bash

sudo docker run --rm -it -v "$PWD":$DOCKER_WORKDIR -w $DOCKER_WORKDIR $DOCKER_IMAGE bash make.bash
rm make.bash
ls -lh $APP-*-*

