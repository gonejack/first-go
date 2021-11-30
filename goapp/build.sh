#!/bin/bash

go version

GOARCH=wasm GOOS=js go build -o web/app.wasm .

go build -o hello .

