#!/usr/bin/env bash

THISDIR=$(dirname $0)

cat <<EOM | gofmt > $THISDIR/bundledSpecInject.go
// Code generated during build process, along with langspec.json. DO NOT EDIT.
package logic

var langSpecJson []byte

func init() {
        langSpecJson = []byte{
        $(cat $THISDIR/langspec.json | hexdump -v -e '1/1 "0x%02X, "' | fmt)
        }
}

EOM