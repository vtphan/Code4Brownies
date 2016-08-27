#!/bin/bash
echo "Compiling c4b_darwin_amd64"
env GOOS=darwin GOARCH=amd64 go build -o c4b_darwin_amd64
echo "Compiling c4b_windows_amd64"
env GOOS=windows GOARCH=amd64 go build -o c4b_windows_amd64.exe
echo "Compiling c4b_linux_amd64"
env GOOS=linux GOARCH=amd64 go build -o c4b_linux_amd64
