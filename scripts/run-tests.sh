#!/bin/bash

set -x

go-bindata -debug -prefix "templates/" templates/

go test -v -race ./...
