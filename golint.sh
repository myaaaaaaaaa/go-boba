#!/bin/bash -ex

go fmt .
go mod tidy
go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

for tape in *.tape; do
	go run github.com/charmbracelet/vhs@latest "$tape"
done

go run ./regen/

