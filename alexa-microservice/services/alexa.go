package services

//// alexa.go:
//// Microservice which commits several queries to convert a spoken question to a spoken answer  
//// Takes a JSON object that has a "speech" field containing .wav(base64) audio data {"speech": "<.wav>"}
//// the sound is converted to text, queried using WolframAlpha, and converted back to speech.
//// Finally the new JSON object containing the .wav(base64) audio data is returned

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

// Recieves a JSON object containing a spoken question in the form of .wav(base64) audio data,
// using the other microservices the data is converted from speech to text, queried using WolframAlpha,
// and converted from text to speech returning an answer to the question in the same format as the input
func wave64QuestionToAnswer(outRsp http.ResponseWriter, inReq *http.Request) {

	if inReq.Method == "POST" {
	
		// get request data
		reqBody, err := ioutil.ReadAll(inReq.Body)
		if err != nil { 
			http.Error(outRsp, "invalid JSON object \nE: " + err.Error(), http.StatusBadRequest)
			return
		}
		inReq.Body.Close()

		// speech to text
		if sttData,	err := queryMicroService(STTPort, STTPath , reqBody); err == nil {

			// WolframAlpha query
			if alphaData, err := queryMicroService(AlphaPort, AlphaPath, sttData); err == nil {

				// text to speech
				if ttsData, err := queryMicroService(TTSPort, TTSPath, alphaData); err == nil {
			
					// return response json speech struct containing .wav(base64) data
					outRsp.WriteHeader(http.StatusOK)
					outRsp.Write(ttsData)
					return
				}
			}
		}
		// reports any errors produced by the microservices
		if err != nil {
			http.Error(outRsp, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(outRsp, "Only POST requests are accepted for Alexa", http.StatusBadRequest)
	}
}

// Sends a request with the input data to the path + port given, then returns the response data
func queryMicroService(localPort string, localPath string, data []byte) ([]byte, error) {

	var status error = nil
	localURI := "http://localhost" + localPort + localPath

	// sends data to other microservice
	client	 := &http.Client{}
	request, err := http.NewRequest("POST", localURI, bytes.NewReader(data))
	if err != nil {
		status = errors.New("Failed to create request for: |" + localURI + "|\nE:" + err.Error())
		return nil, status
	}

	// commit microservice query
	response, err := client.Do(request)
	if err != nil { 
		status = errors.New("|" + localURI + "| failed to respond" + "\nE:" + err.Error())
		return nil, status
	}
	
	// get body data
	rspBody, err := ioutil.ReadAll(response.Body)
	if err != nil { 
		status = errors.New("Failed to read response data from: |" + localURI + "|\nE:" + err.Error())
		return nil, status
	}
	response.Body.Close()

	return rspBody, status
}

// Sets a listener thread for path "/alexa" on port ":3000" 
func SetAlexaListenerThread() {

	router := mux.NewRouter()
	router.HandleFunc(AlexaPath, wave64QuestionToAnswer).Methods("POST")

	// set listen to wait for request
	if err := http.ListenAndServe(AlexaPort, router); err != nil {
		panic(err)
	}
}