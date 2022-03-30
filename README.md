# Alexa-Microservice
Runs a local server hosting four microservices for: Text to Speech, Speech to Text, Text Question to Text Answer, Spoken Question to Spoken Answer.
  
## Content
 - [How To Run?](#How-To-Run?)
   - [Requirements](#Requirements)
 - [Project Guide](#Project-Guide)
    
## How To Run?
 * See Requirements!
 1. Open `two` terminals and follow step `2` for both of them
 2. Navigate to the `alexa-microservice` directory containing `alexa-microservice.go` (and the `run-server` scripts)
 3. In the `first terminal` use command: `bash run-server.sh` for Linux, or `.\run-server.ps1` for Windows Powershell
 4. In the `second terminal` KEEP `alexa-microservice` as the root directory, and from here run the `test` scripts
    * Example: `bash test/<name-test>.sh`
  
**Issues**
 * If step `3` failed, ensure the correct permissions for all files within the `alexa-microservice` directory are granted
 * `base64` can vary across `Linux` versions, which can cause incompatable whitespace being encoded . . .
   * For Ubuntu: `base64 -i` should work fine
   * If this command fails, use: `base64 -w 0 -i`
   * Finally if both fail, the input is most likely incorrect
  
### Requirements
 * `gorilla mux` must be install, use command: `go get -u github.com/gorilla/mux`
 * Two keys must be provided for `WolframAlpha` and `Microsoft Cognitive Services` inside the `config.json` file
   * If you do not have keys for these services they can be obtained via their respective sites
  
## Project Guide
All four microservices expect a specific JSON object as input.  
 * `alpha.go` and `tts.go` require a JsonText object which contains a text query
   * Example: `{"text": "What is the melting point of silver?"}`
 * `stt.go` and `alexa.go` require a JsonSpeech object which contains a spoken query as a base64 encoded `.wav` (WAVE) file
   * Example: `{"speech": "<.wav>(base64)"}`
  
| Directories | |
|--:                    |:-- |
|`bin/`                 | Hosts the exeutable and any files it needs |
|`config/`              | Contains `config.go` a handler for `config.json` which stores vital information for all the microservices to run |
|`res/`                 | Holds various resource files for the user to utilise or to be used as a file output dump during the programs debug method |
|`services/`            | Package directory storing the microservices used in the program |
|`test/`                | Stores all the test scripts with both bash and powershell variants |
  
| Root Files | |
|--:                    |:-- |
|`alexa-microservice.go`| The `main` program, runs several go routines to create multiple listerner threads for the microservices |
|`go.mod`| Hosts the `alexa-microservice` module requirements |
|`run-server.sh`| Script to build and run the `alexa-microservice server` (for Bash) |
|`run-server.ps1`| Windows Powershell variant of `run-server.sh` |

| Services  | |
|--:        |:-- |
|`alpha.go` | Microservice which provides computational knowledge via communication with WolframAlpha's API |
|`stt.go`   | Microservice which computes `Speech to Text` using Microsoft Cognitive Services (MCS) |
|`tts.go`   | Microservice which computes `Text to Speech` using Microsoft Cognitive Services (MCS) |
|`alexa.go` | Microservice which commits several queries to convert a spoken question to a spoken answer |
|`utils.go` | Hosts utility functions, variables, & structs used in data representation |
  
**Debug: Write Resource to File**  
In `services/utils.go` if `debugWriteResourceToFile` is set to `true`,  
files managed by each microservice will be written to the resource directory (`/res`)  