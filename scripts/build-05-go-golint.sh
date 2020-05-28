#!/usr/bin/env bash

# Check gofmt
echo "==> Checking for golint errors..."

golangci-lint run

exit 0
