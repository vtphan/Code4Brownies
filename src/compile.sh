#!/bin/bash
version=`cat VERSION`
echo "Compile c4b binaries for version $version."
echo "Compiling c4b_osx_$version"
env GOOS=darwin GOARCH=amd64 go build -o c4b_osx_$version

echo "Compiling c4b_win_$version.exe"
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o c4b_win_$version.exe
