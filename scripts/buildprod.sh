#!/bin/bash

# Build static linux/amd64 binary to match Dockerfile base image
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notely .
