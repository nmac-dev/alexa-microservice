#!/bin/sh
echo     '{"text": "What is the melting point of silver?"}' > bin/tts-input.json
TTS_OUT=`curl -s -X POST -d "@bin/tts-input.json" localhost:3003/tts`

echo $TTS_OUT
