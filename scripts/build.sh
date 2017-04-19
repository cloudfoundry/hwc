#!/usr/bin/env bash

set -ex

mkdir -p $PWD/hwc-rel
GOOS=windows go build -o $PWD/hwc-rel/hwc.exe code.cloudfoundry.org/hwc
