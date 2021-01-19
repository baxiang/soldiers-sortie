#!/bin/bash
mkdir pkg/pb
protoc -I ./api --go_out=plugins=grpc:./pkg/pb/ api/*.proto
CGO_ENABLED=0  go build -o staticservice cmd/staticservice/main.go
