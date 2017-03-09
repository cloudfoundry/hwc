#!/usr/bin/env bash

ROOTDIR="$( dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )" )"
BINDIR=$ROOTDIR/bin

mkdir -p $BINDIR

set -ex

GOOS=windows go build -o $BINDIR/hwc.exe github.com/cloudfoundry-incubator/hwc/hwc
