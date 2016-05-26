#!/bin/bash
env GOOS=darwin GOARCH=amd64 go build -o c4b_darwin_amd64
env GOOS=windows GOARCH=amd64 go build -o c4b_windows_amd64.exe
