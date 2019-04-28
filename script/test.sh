#!/bin/sh
# Run package tests for a file/directory, or all tests if no argument is passed.
# Useful to e.g. execute package tests for the file currently open in Vim.
# Usage: script/test.sh [path]

export PATH="$PATH:$GOPATH/bin"

set -e

go_pkg_from_path() {
    path=$1
    if test -d "$path"; then
        dir="$path"
    else
        dir=$(dirname "$path")
    fi
    (cd "$dir" && go list)
}

if test $# -gt 0; then
    pkg=$(go_pkg_from_path "$1")
    verbose=-v
else
    pkg=./...
    verbose=
fi

exec go test ${GOTESTOPTS:-$verbose} "$pkg"