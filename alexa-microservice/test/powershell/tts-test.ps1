#!/bin/ps1
$JSON_INPUT='{"text": "What is the melting point of silver?"}'
$TTS_OUT = Invoke-RestMethod -Method POST -Body $JSON_INPUT -Uri @localhost:3003/tts `
            | ConvertTo-Json -Compress

Write-Output $TTS_OUT