#!/bin/ps1
$WAVE = [convert]::ToBase64String((Get-Content -Path ".\res\tts-speech.wav" -Encoding Byte))
$JSON='{"speech": "' + $WAVE + '"}'
$JSON2 = Invoke-RestMethod -Method POST -Body $JSON -Uri @localhost:3002/stt | ConvertTo-Json -Compress

Write-Output $JSON2


# JSON="{\"speech\":\"‘base64 -i speech.wav‘\"}"
# JSON2=‘curl -s -X POST -d "$JSON" localhost:3002/stt‘
# echo $JSON2