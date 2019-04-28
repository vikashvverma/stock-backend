#!/bin/sh
# Check code style and correctness.
# Usage: script/lint.sh

export PATH="$PATH:$GOPATH/bin"

golint $(go list ./... | grep -v /vendor/)
go vet $(go list ./... | grep -v /vendor/)