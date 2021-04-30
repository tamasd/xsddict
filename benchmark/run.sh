#!/bin/bash

go run ../cmd/main.go -extra=data/extra.txt data/schema.xsd dict.dat

go run compress.go
