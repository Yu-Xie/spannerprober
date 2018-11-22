#!/usr/bin/env bash

goimports -w $(find . -type f -name '*.go' -not -path "./vendor/*")
gofmt -w $(find . -type f -name '*.go' -not -path "./vendor/*")
