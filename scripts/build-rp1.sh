#!/bin/sh

source ./build-common.sh
source ./deployable.sh

export GOOS=linux
export GOARCH=arm
export GOARM=5

build_app

create_deployable
