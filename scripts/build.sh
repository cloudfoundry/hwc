#!/usr/bin/env bash

set -ex

mkdir -p $PWD/hwc-rel
# x64
CGO_ENABLED=1 GO_EXTLINK_ENABLED=1 CC="x86_64-w64-mingw32-gcc" GOOS="windows" GOARCH="amd64" go build -o $PWD/hwc-rel/hwc.exe code.cloudfoundry.org/hwc
# Win32
CGO_ENABLED=1 GO_EXTLINK_ENABLED=1 CC="i686-w64-mingw32-gcc" GOOS="windows" GOARCH="386" go build -o $PWD/hwc-rel/hwc_x86.exe code.cloudfoundry.org/hwc