package services

//// tts.go: (
//// Microservice which computes Text to Speech using Microsoft Cognitive Services (MCS).
//// Takes a JSON object containing text to be transformed to speech {"text": "<TextToSpeech>"}
//// and returns a JSON object containing a WAVE(base64) file's speech data {"speech": "<.wav>"}

import (
	"alexa-microservice/config"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

var (
	// XML header with config values
	xmlHeader = `<?xml version="`+ config.GetXmlVersion() +
				`" encoding="`	 + config.GetXmlEncoding() +`"?>` + "\n"

	// SSML Tag: <speak> <voice> </speak>
	ssmlSpeak = SSMLSpeak {
		XMLName:	xml.Name{},
		Version:	config.GetXmlVersion(),
		Lang:		config.GetXmlLang(),
		Voice: SSMLVoice{
			XMLName:	xml.Name{},
			Lang:		config.GetXmlLang(),
			Name:		config.GetSsmlVoiceName(),
		},
	}
	// sets input as ssml & xml
	hTypeTTS = rHeader{Key: "Content-Type",				Value: "application/ssml+xml"		}
	// sets output as RIFF|WAVE(base64) {.wav}
	hOutTTS	 = rHeader{Key: "X-Microsoft-OutputFormat", Value: "riff-16khz-16bit-mono-pcm"	}

	// MicrosoftCS TTS URI
	ttsURI = "https://"+ config.GetMcsRegion() + config.GetMcsTTS() + config.GetMcsSrvPath()
)

// Decodes text from inbound JSON request, then encodes SSML for Microsoft Cognitive Services
// and send an output response containing the .wav(base64) speech data
func textToSpeech(outRsp http.ResponseWriter, inReq *http.Request) {

	if inReq.Method == "POST" {

		// json data -> text struct
		text := JsonText{Data: ""}
		err  := json.NewDecoder(inReq.Body).Decode(&text)
		
		if err != nil || text.Data == "" {
			http.Error(outRsp, "invalid JSON object", http.StatusBadRequest)
			return
		}

		// json request -> ssml xml
		ssmlSpeak.Voice.CharData = []byte(text.Data)

		xmlElements, err := xml.MarshalIndent(ssmlSpeak, "", "\t")
		if err != nil {
			http.Error(outRsp, "failed to generate XML elements", http.StatusBadRequest)
			return
		}
		xmlElements = []byte(xmlHeader + string(xmlElements))

		// query MCS
		if response, err := ttsCommit(xmlElements); err == nil {
			
			defer response.Body.Close()
			
			body, err := ioutil.ReadAll(response.Body)
			if err != nil { 
				http.Error(outRsp, "Failed to read Microsoft Cognitive Services", http.StatusBadRequest)
				return
			}
			outRsp.WriteHeader(http.StatusOK)

			//// Debug Only
			if debugWriteResourceToFile {	
				writeToResFile("tts-ssml.xml", xmlElements)
				writeToResFile("tts-speech.wav", body)
			}

			// WAVE(base64) data -> json speech struct
			json.NewEncoder(outRsp).Encode(JsonSpeech{Data: body})

		} else {
			http.Error(outRsp, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(outRsp, "Only POST requests are accepted for Text to Speech", http.StatusBadRequest)
	}
}

// Sends a request with SSML to Microsoft Cognitive Services, then returns the response
func ttsCommit(xmlElements []byte) (*http.Response, error) {

	var status error = nil

	// sends XML to MCS
	client	 := &http.Client{}
	request, err := http.NewRequest("POST", ttsURI, bytes.NewReader(xmlElements))
	if err != nil {
		status = errors.New("Failed to create request for: " + ttsURI + "\nE:" + err.Error())
	}

	request.Header.Set(hTypeTTS.Key, 	hTypeTTS.Value)
	request.Header.Set(hAzureKey.Key, 	hAzureKey.Value)
	request.Header.Set(hOutTTS.Key, 	hOutTTS.Value)

	// commit tts query
	response, err := client.Do(request)
	if err != nil { 
		status = errors.New("Microsoft Cognitive Services failed to respond" + "\nE:" + err.Error())
	}

	return response, status
}

// Sets a listener thread for path "/tts" on port ":3003" 
func SetTTSListenerThread() {

	router := mux.NewRouter()
	router.HandleFunc(TTSPath, textToSpeech).Methods("POST")

	// set listen to wait for request
	if err := http.ListenAndServe(TTSPort, router); err != nil {
		panic(err)
	}
}