#!/bin/bash -e

go test ./...

~/go/bin/wire cmd/gobot/main.go
go build -o gobot cmd/gobot/wire_gen.go
rm -f cmd/gobot/wire_gen.go
