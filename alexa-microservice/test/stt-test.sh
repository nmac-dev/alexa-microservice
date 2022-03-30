#!/bin/sh
WAVE64=`base64 -w 0 -i ./res/tts-speech.wav`
echo    '{"speech": "'$WAVE64'"}' > bin/stt-input.json
STT_OUT=`curl -s -X POST -d "@bin/stt-input.json" localhost:3002/stt`

echo $STT_OUT
