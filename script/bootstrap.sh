#!/bin/sh
# Install Go dependencies.
# Usage: script/bootstrap.sh

export PATH="$PATH:$GOPATH/bin"

set -e

go get -d ./...
go get  \
    golang.org/x/lint/golint \
    github.com/mitchellh/gox
go mod vendor