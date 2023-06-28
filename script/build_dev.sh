#!/bin/bash

export GOOS=linux
export GOARCH=amd64
go build -o build/capuchin_linux_amd64 -trimpath -ldflags "-s -w"
