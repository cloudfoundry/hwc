#!/bin/bash

set -eu
set -o pipefail

function run() {
    local arch os output
    arch="${1:?Please provide an architecture}"
    os="${2:?Please provide an OS}"
    version="${3:?Please provide a version}"
    output="${4:?Please provide an output directory}"

    local binary tarball
    binary="hwc-${os}-${arch}"
    tarball="${binary}.tgz"

    local cc
    if [[ "$arch" == "amd64" ]]; then
        cc="x86_64-w64-mingw32-gcc"
    else
        cc="i686-w64-mingw32-gcc"
    fi

    CGO_ENABLED=1 GO_EXTLINK_ENABLED=1 CC="$cc" GOARCH="${arch}" GOOS="${os}" \
    go build \
    -o "${output}/${binary}" \
    -ldflags "-X main.version=${version}"

    tar -czvf "${output}/${tarball}" -C "${output}" "${binary}"
}

run "$@"
