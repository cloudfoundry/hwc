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

    GOARCH="${arch}" GOOS="${os}" \
    go build \
    -o "${output}/${binary}" \
    -ldflags "-X main.version=${version}"

    tar -czvf "${output}/${tarball}" -C "${output}" "${binary}"
}

run "$@"
