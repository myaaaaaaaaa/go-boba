#!/bin/bash -ex

go get -u ./...
rm go.sum
sed -i -e 's/.*indirect.*//' go.mod
go mod tidy
