#!/bin/bash

go mod download && go mod verify && go mod tidy

for GOOS in windows darwin linux; do
  for GOARCH in 386 amd64; do
     export GOOS GOARCH
     go build -v -o ./build/bin/server-$GOOS-$GOARCH ./cmd/server
     go build -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -v -o ./build/bin/client-$GOOS-$GOARCH ./cmd/client
  done
done
