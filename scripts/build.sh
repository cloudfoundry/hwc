#!/usr/bin/env bash

set -ex

mkdir -p $PWD/hwc-rel
CGO_ENABLED=1 GO_EXTLINK_ENABLED=1 CC="x86_64-w64-mingw32-gcc" GOOS=windows go build -o $PWD/hwc-rel/hwc.exe code.cloudfoundry.org/hwc
