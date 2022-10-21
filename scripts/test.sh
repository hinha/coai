#!/bin/bash
set -e

readonly service="$1"

cd "./core/$service"
# shellcheck disable=SC2046
go test -count=1 -p=8 -parallel=8 -race ./...
