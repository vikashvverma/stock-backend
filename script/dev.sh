#!/bin/bash
# Run project for current architecture.
# Usage: script/dev.sh
export PATH="$PATH:$GOPATH/bin"
export config=config/config.json

if [[ $(uname -s) == "Darwin" ]]; then
  ext=""
  platform="darwin"
elif [[ $(uname -s) == *"MINGW"* ]]; then
  ext=".exe"
  platform="windows"
  tag=""
else
  ext=""
  platform="linux"
fi

arch=amd64
svc="stock"
version=0.1
build=`date +%FT%T%z`

build_dir="build/latest"

echo "--> Building binary for ${platform}"
gox \
    -output="${build_dir}/${svc}_${version}_{{.OS}}_{{.Arch}}/${svc}" \
    -os="${platform}" \
    -arch="amd64" \
    -ldflags "-X main.version=$version -X main.build=$build" \
    ./cmd/stock

program=./build/latest/${svc}_"${version}"_"${platform}_${arch}"/"${svc}${ext}"

echo "--> running application..."

exec "${program}" \
     -config ${config}