#!/bin/ps1
$JSON_INPUT='{"text": "What is the melting point of silver?"}'
$ALPHA_OUT = Invoke-RestMethod -Method POST -Body $JSON_INPUT -Uri @localhost:3001/alpha `
            | ConvertTo-Json -Compress

Write-Output $ALPHA_OUT