#!/usr/bin/env bash

set -ex

GOOS=windows go build -o hwc-rel/hwc.exe github.com/cloudfoundry-incubator/hwc
