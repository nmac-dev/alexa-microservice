#!/bin/ps1

$OUT_DIR = ".\bin\"
$TARGET = $OUT_DIR + "alexa-microservice"

# create bin directory
if ( !(Test-Path $OUT_DIR) ) {
    New-Item -itemType Directory -Path $OUT_DIR
}
# build
go mod tidy
go build -o $OUT_DIR

# run server
Write-Output "Server Running: use 'Ctrl + C' to terminate. . ."
Invoke-Expression $TARGET