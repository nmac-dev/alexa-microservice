package services

//// alpha.go: 
//// Microservice which provides computational knowledge via communication with WolframAlpha's API.
//// Takes a JSON object containing a text based question {"text": "<question>"}, 
//// and returns a JSON object containing a text based answer for the question {"text": "<answer>"}

import (
	"alexa-microservice/config"
	"net/http"
	"encoding/json"
	"errors"
	"io/ioutil"
	"github.com/gorilla/mux"
)

var (
	alphaQueryKey = "i"
	alphaQueryURI = config.GetAlphaApiURI() + config.GetAlphaAppPath() + config.GetAlphaAppID()
)

// Recieves a JSON request, it is decoded and send as a query to WolframAlpha,
// then the response from WolframAlpha is encoded back to JSON and sent to the requestee
func alphaQuery(outRsp http.ResponseWriter, inReq *http.Request) {

	if inReq.Method == "POST" {

		// parse json data -> text struct
		text := JsonText{Data: ""} 
		err  := json.NewDecoder(inReq.Body).Decode(&text)

		// json parse failed or text is invalid query
		if err != nil || text.Data == "" {
			http.Error(outRsp, "invalid JSON object", http.StatusBadRequest)
			return
		}

		// encode response to json
		if response, err := alphaCommit(text.Data); err == nil {

			// ensures IO for response closes on stack call
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil { 
				http.Error(outRsp, "Failed to read JSON from WolframAlpha", http.StatusBadRequest)
			}
			outRsp.WriteHeader(http.StatusOK)
			
			text.Data = string(body) 

			//// Debug Only
			if debugWriteResourceToFile {
				writeToResFile("alpha-answer.json", *singular(json.Marshal(text)))
			}

			// encode text struct to json
			json.NewEncoder(outRsp).Encode(text)

		} else {
			http.Error(outRsp, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(outRsp, "Only POST requests are accepted for Alpha Queries", http.StatusBadRequest)
	}
}

// Builds a new request query, then commits it to the URI and catches the response to be returned
func alphaCommit(text string) (*http.Response, error) {

	var status error = nil

	// sends JSON object to WolframAlpha API
	client	 := &http.Client{}
	request, err := http.NewRequest("POST", alphaQueryURI, nil)
	if err != nil {
		status = errors.New("Failed to create request for: " + alphaQueryURI + "\nE:" + err.Error())
	}

	// append text query to URI
	query := request.URL.Query()
	query.Add(alphaQueryKey, text)
	request.URL.RawQuery = query.Encode()

	// commit alpha query
	response, err := client.Do(request)
	if err != nil { 
		status = errors.New("WolframAlpha query failed to respond" + "\nE:" + err.Error())
	}

	return response, status
}

// Sets a listener thread for path "/alpha" on port ":3001" 
func SetAlphaListenerThread() {

	router := mux.NewRouter()
	router.HandleFunc(AlphaPath, alphaQuery).Methods("POST")

	// set listen to wait for request
	if err := http.ListenAndServe(AlphaPort, router); err != nil {
		panic(err)
	}
}