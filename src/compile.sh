#!/bin/bash
version=`cat VERSION`
echo "Compile c4b binaries for version $version."
echo "Compiling c4b_darwin_$version"
env GOOS=darwin GOARCH=amd64 go build -o c4b_darwin_$version

# echo "Compiling c4b_windows_$version"
# env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -o c4b_windows_$version.exe
