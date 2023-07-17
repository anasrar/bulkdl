#!/usr/bin/env sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-s -w -X 'main.VERSION=test'" -o bin/bulkdl_linux_amd64 cmd/bulkdl/*.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build --ldflags="-s -w -X 'main.VERSION=test'" -o bin/bulkdl_windows_amd64.exe cmd/bulkdl/*.go
