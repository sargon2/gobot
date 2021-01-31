#!/bin/bash

go test ./...
go build -o gobot cmd/gobot/main.go
