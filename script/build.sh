#!/bin/sh
# Build project for multiple architectures. With --release, also create
# download archives and other release artifacts.
# Usage: script/build [-r|--release]
export PATH="$PATH:$GOPATH/bin"

set -e

git_version=0.1
version=$(expr "$git_version" : v*'\(.*\)' | sed -e 's/-/./g')
build_dir="build/$version"

echo "Building stock $git_version ..."

echo "Installing dependencies..."
script/bootstrap.sh

echo "Running lint checks..."
script/lint.sh

echo "Running tests..."
script/test.sh

echo "Cross-compiling binaries..."
rm -rf "$build_dir"
gox \
    -output="${build_dir}/{{.Dir}}_${version}_{{.OS}}_{{.Arch}}/{{.Dir}}" \
    -os="darwin linux windows freebsd openbsd" \
    -arch="amd64" \
    -ldflags "-X main.version=$git_version" \
    ./...

rm -rf build/latest
ln -snf "$version" build/latest

case "$1" in
-r|--release)
    echo "Creating zip archives..."
    cd "$build_dir"
    for i in *; do zip -r "$i.zip" "$i"; done

    echo "Creating SHA256SUMS file..."
    if hash shasum 2>/dev/null; then
        shasum -a256 *.zip > SHA256SUMS
    else
        sha256sum *.zip > SHA256SUMS
    fi

esac

echo "Done."