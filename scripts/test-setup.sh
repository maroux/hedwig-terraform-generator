#!/bin/bash

set -x

if [[ "${TRAVIS}" == "true" ]]; then
    go get github.com/kardianos/govendor

    go tool vet .

    gofmt -d .
    if [ "$?" -ne "0" ]; then
        echo "Code not properly formatted. Run 'go fmt'"
        exit 1
    fi

    govendor sync

    if [ ! -z "$(govendor list -no-status +outside)" ]; then
        echo "External dependencies found, only vendored dependencies supported"
        exit 1
    fi

    # XXX because govendor sync doesn't install things in BIN folder
    go get github.com/go-bindata/go-bindata/...
fi
