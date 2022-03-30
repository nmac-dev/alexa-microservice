#!/bin/sh
echo     '{"text": "What is the melting point of silver?"}' > JSON_INPUT
TTS_OUT=`curl -s -X POST -d "@JSON_INPUT" localhost:3003/tts`

echo $TTS_OUT
