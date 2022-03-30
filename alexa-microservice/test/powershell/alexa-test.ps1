#!/bin/ps1
$WAVE64  = [convert]::ToBase64String((Get-Content -Path ".\res\alexa-question.wav" -Encoding Byte))
$JSON_INPUT  = '{"speech": "' + $WAVE64 + '"}'
$ALEXA_OUT = Invoke-RestMethod -Method POST -Body $JSON_INPUT -Uri @localhost:3000/alexa `
            | ConvertTo-Json -Compress

# write .wav(base64) data from speech JSON object to file
$WAVE64 = [convert]::FromBase64String(( (($ALEXA_OUT) -Split '"')[3] ))
Set-Content -Path ".\res\alexa-answer.wav" -Value $WAVE64  -Encoding Byte