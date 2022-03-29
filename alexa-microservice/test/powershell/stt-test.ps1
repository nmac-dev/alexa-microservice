#!/bin/ps1
$WAVE = [convert]::ToBase64String((Get-Content -Path ".\res\tts-speech.wav" -Encoding Byte))
$JSON='{"speech": "' + $WAVE + '"}'
$JSON2 = Invoke-RestMethod -Method POST -Body $JSON -Uri @localhost:3002/stt | ConvertTo-Json -Compress

Write-Output $JSON2