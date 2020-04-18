#!/bin/bash

rm -r ./bin
mkdir ./bin

go build -o ./bin/app cmd/main.go

