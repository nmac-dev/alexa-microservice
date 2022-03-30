#!/bin/ps1
$WAVE64 = [convert]::ToBase64String((Get-Content -Path ".\res\tts-speech.wav" -Encoding Byte))
$JSON_INPUT='{"speech": "' + $WAVE64 + '"}'
$STT_OUT = Invoke-RestMethod -Method POST -Body $JSON_INPUT -Uri @localhost:3002/stt `
            | ConvertTo-Json -Compress

Write-Output $STT_OUT