#!/bin/sh
echo    '{"text": "What is the melting point of silver?"}' > JSON_INPUT
ALPHA_OUT=`curl -s -X POST -d "@JSON_INPUT" localhost:3001/alpha`

echo $ALPHA_OUT
