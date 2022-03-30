#!/bin/sh
WAVE64=`base64 -w 0 -i ./res/alexa-question.wav`
echo    '{"speech": "'$WAVE64'"}' > bin/alexa-input.json
ALEXA_OUT=`curl -s -X POST -d "@bin/alexa-input.json" localhost:3000/alexa`

echo $ALEXA_OUT | cut -d '"' -f4 | base64 -w 0 -d > ./res/alexa-answer.wav
