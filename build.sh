#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o ./builds/yer_linux_amd64
GOOS=linux GOARCH=arm64 go build -o ./builds/yer_linux_arm64
GOOS=darwin GOARCH=amd64 go build -o ./builds/yer_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -o ./builds/yer_darwin_arm64
GOOS=windows GOARCH=amd64 go build -o ./builds/yer_windows_amd64.exe

cp ./builds/yer* ~/dev/yer_app/public/scripts