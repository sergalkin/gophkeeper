#!/bin/bash

go mod download && go mod verify

CompileDaemon --build="go build -v -o ./.tmp ./cmd/..." --command=./.tmp/server