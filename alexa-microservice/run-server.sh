#!/bin/sh

OUT_DIR="./bin/"
TARGET=$OUT_DIR"alexa-microservice"

# create bin directory
[[ -d $OUT_DIR ]] || mkdir $OUT_DIR

# build
go mod tidy
go build -o $OUT_DIR

# run server
echo "Server Running: use 'Ctrl + C' to terminate. . ."
$TARGET