#!/bin/bash -ex

go fmt .
go mod tidy
# go tool modernize -fix -test ./...

for tape in *.tape; do
	go tool vhs "$tape"
done

go run ./regen/

