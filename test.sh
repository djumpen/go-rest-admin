#!/bin/bash

### FORMAT
gofmt -s -w -l .
go vet github.com/djumpen/go-rest-admin
go get -u github.com/alecthomas/gometalinter
gometalinter ./... --exclude="\bexported \w+ (\S*['.]*)([a-zA-Z'.*]*) should have comment or be unexported\b" --install

### TEST (enable it after adding tests)
go test ./...

