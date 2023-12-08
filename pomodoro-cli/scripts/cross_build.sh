#!/bin/bash

OS_LIST="linux windows darwin"
ARCH_LIST="amd64 arm arm64"

for os in ${OS_LIST}; do
  for arch in ${ARCH_LIST}; do
    if [[ "$os/$arch" =~ ^(windows/arm64|darwin/arm)$ ]]; then continue; fi
    echo "Building binary for $os $arch"
    mkdir -p releases/${os}/${arch}
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -tags=inmemory -o releases/${os}/${arch}/
  done
done