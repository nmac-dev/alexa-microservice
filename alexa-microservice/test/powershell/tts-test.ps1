#!/bin/ps1
$JSON='{"text": "What is the melting point of silver?"}'
$JSON2 = Invoke-RestMethod -Method POST -Body $JSON -Uri @localhost:3003/tts | ConvertTo-Json -Compress

Write-Output $JSON2