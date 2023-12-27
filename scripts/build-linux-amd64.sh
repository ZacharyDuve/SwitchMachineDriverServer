#!/bin/sh

source ./build-common.sh
source ./deployable.sh

export GOOS=linux
export GOARCH=amd64

build_app

create_deployable