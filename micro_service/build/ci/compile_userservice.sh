#!/usr/bin/env bash
mkdir pkg/pb
protoc -I ./api --go_out=plugins=grpc:./pkg/pb/ api/*.proto
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags netgo -ldflags="-s -w" -tags netgo -o userservice cmd/userservice/main.go