#!/bin/ps1
$WAVE  = [convert]::ToBase64String((Get-Content -Path ".\res\alexa-question.wav" -Encoding Byte))
$JSON  = '{"speech": "' + $WAVE + '"}'
$JSON2 = Invoke-RestMethod -Method POST -Body $JSON -Uri @localhost:3000/alexa | ConvertTo-Json -Compress

$WAVE2 = [convert]::FromBase64String(( (($JSON2) -Split '"')[3] ))
# write .wav(base64) data from speech JSON object to file
Set-Content -Path ".\res\alexa-answer.wav" -Value $WAVE2  -Encoding Byte