#!/bin/bash -ex

go fmt .
go mod tidy
go tool modernize -fix -test ./...
