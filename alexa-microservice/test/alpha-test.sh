#!/bin/sh
echo    '{"text": "What is the melting point of silver?"}' > bin/alpha-input.json
ALPHA_OUT=`curl -s -X POST -d "@bin/alpha-input.json" localhost:3001/alpha`

echo $ALPHA_OUT
