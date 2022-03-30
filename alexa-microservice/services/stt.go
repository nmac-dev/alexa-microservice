package services

//// stt.go:
//// Microservice which computes Speech to Text using Microsoft Cognitive Services (MCS).
//// Takes a JSON object that has a "speech" field contaning .wav(base64) sound data {"speech": "<.wav>"}
//// the sound from the data is transformed to text as a JSON object {"text": "<TextToSpeech>"}
//// finally the new JSON object is returned 

import (
	"alexa-microservice/config"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

var (
	// sets input RIFF|WAVE(base64)
	hTypeWAVE64 = rHeader{Key: "Content-Type", Value: "audio/wav;codecs=audio/pcm;samplerate=16000"}
	
	// MicrosoftCS STT URI
	sttURI = "https://" + 
			config.GetMcsRegion()+ config.GetMcsSTT()+ config.GetMcsSrvPath()+ config.GetMcsSTTLang()
)

// Decodes .wav(base64) data from inbound JSON object, then is given to Microsoft Cognitive Services
// and the response is returned as a text JSON object which contains the speech to text conversion   
func speechToText(outRsp http.ResponseWriter, inReq *http.Request) {

	if inReq.Method == "POST" {

		// json data -> speech struct
		speechWAVE := JsonSpeech{}
		err  := json.NewDecoder(inReq.Body).Decode(&speechWAVE)
		
		if err != nil || speechWAVE.Data == nil {
			http.Error(outRsp, "invalid JSON object \nE: " + err.Error(), http.StatusBadRequest)
			return
		}
		
		// query MCS
		if response, err := sttCommit(speechWAVE.Data); err == nil {
			
			defer response.Body.Close()
			
			body, err := ioutil.ReadAll(response.Body)
			if err != nil { 
				http.Error(outRsp, "Failed to read Microsoft Cognitive Services", http.StatusBadRequest)
				return
			}
			outRsp.WriteHeader(http.StatusOK)

			// MCS json -> json text struct
			mcsJson := MCSJsonRsp{}
			json.Unmarshal(body, &mcsJson)

			speechText := JsonText{Data: mcsJson.DisplayText}

			//// Debug Only
			if debugWriteResourceToFile {	
				writeToResFile("stt-text.json", *singular(json.Marshal(speechText)))
			}

			// return json text struct
			json.NewEncoder(outRsp).Encode(speechText)

		} else {
			http.Error(outRsp, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(outRsp, "Only POST requests are accepted for Speech to Text", http.StatusBadRequest)
	}
}

// Sends a request with .wav(base64) data to Microsoft Cognitive Services, then returns the response
func sttCommit(speechData []byte) (*http.Response, error) {

	var status error = nil

	// sends WAVE64 to MCS
	client	 := &http.Client{}
	request, err := http.NewRequest("POST", sttURI, bytes.NewReader(speechData))
	if err != nil {
		status = errors.New("Failed to create request for: " + sttURI + "\nE:" + err.Error())
	}

	request.Header.Set(hTypeWAVE64.Key,	hTypeWAVE64.Value)
	request.Header.Set(hAzureKey.Key,	hAzureKey.Value)

	// commit stt query
	response, err := client.Do(request)
	if err != nil { 
		status = errors.New("Microsoft Cognitive Services failed to respond" + "\nE:" + err.Error())
	}

	return response, status
}

// Sets a listener thread for path "/stt" on port ":3002" 
func SetSTTListenerThread() {

	router := mux.NewRouter()
	router.HandleFunc(STTPath, speechToText).Methods("POST")

	// set listen to wait for request
	if err := http.ListenAndServe(STTPort, router); err != nil {
		panic(err)
	}
}