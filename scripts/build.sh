#!/bin/bash -e

cd "$(dirname "${BASH_SOURCE[0]}")"
cd ..

go test ./...

~/go/bin/wire cmd/gobot/main.go
go build -o gobot cmd/gobot/wire_gen.go
rm -f cmd/gobot/wire_gen.go
