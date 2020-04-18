#!/bin/bash
#COMMAND FOR DOCKER ONLY.
##To run the app outside of the go path uncomment.
# go get -u github.com/golang/dep/cmd/dep
#dep init
# dep ensure
#gofmt -s -w -l .
rm -r ./bin
mkdir ./bin
go build -o ./bin/app cmd/*.go
bin/./app
