#!/bin/sh
WAVE64=`base64 -w 0 -i ./res/tts-speech.wav`
echo    '{"speech": "'$WAVE64'"}' > JSON_INPUT
STT_OUT=`curl -s -X POST -d "@JSON_INPUT" localhost:3002/stt`

echo $STT_OUT
