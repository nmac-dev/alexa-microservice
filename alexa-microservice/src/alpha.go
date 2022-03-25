package src

import (
	"alexa-microservice/config"
	"net/http"
	"encoding/json"
	"errors"
	"io/ioutil"
	"github.com/gorilla/mux"
)

const (
	queryKey  = "i"
	alphaPath = "/alpha"
	alphaPort = ":3001"
)

// Recieves a JSON request, it is decoded and send as a query to WolframAlpha,
// then the response from WolframAlpha is encoded back to JSON and sent to the requestee
func alphaQuery(outRsp http.ResponseWriter, inReq *http.Request) {

	// ignore GET requests
	if inReq.Method == "POST" {

		// parse json data -> text struct
		text := Text{Data: ""} 
		err  := json.NewDecoder(inReq.Body).Decode(&text)

		// json parse failed or text is invalid query
		if err != nil || text.Data == "" {
			http.Error(outRsp, "invalid JSON object", http.StatusBadRequest)
			return
		}

		// encode response to json
		if response, err := commitQuery(inReq.Method, text.Data); err == nil {

			// ensures IO for response closes on stack call
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil { 
				http.Error(outRsp, "Failed to read JSON from WolframAlpha", http.StatusBadRequest)
			}

			outRsp.WriteHeader(http.StatusOK)

			// encode text struct to json
			text.Data = string(body) 
			json.NewEncoder(outRsp).Encode(text)

		} else {
			http.Error(outRsp, err.Error(), http.StatusBadRequest)
		}
	}
}

// Builds a new request query, then commits it to the URI and catches the response to be returned
func commitQuery(method string, text string) (*http.Response, error) {

	var status error = nil

	// empty client struct
	client	 := &http.Client{}
	alphaURI := config.GetAlphaApiURI() + config.GetAlphaAppPath() + config.GetAlphaAppID()
	
	// build request
	request, err := http.NewRequest(method, alphaURI, nil)
	if err != nil {
		status = errors.New("Failed to create request for: " + alphaURI + "\nE:" + err.Error())
	}

	// build query
	query := request.URL.Query()
	query.Add(queryKey, text)
	request.URL.RawQuery = query.Encode()

	// commit query
	response, err := client.Do(request)
	if err != nil { 
		status = errors.New("WolframAlpha query failed to respond" + "\nE:" + err.Error())
	}

	return response, status
}

// Sets a listener thread for path "/alpha" on port ":3001" 
func SetAlphaListenerThread() {

	router := mux.NewRouter()
	router.HandleFunc(alphaPath, alphaQuery).Methods("POST")

	// set listen to wait for request
	if err := http.ListenAndServe(alphaPort, router); err != nil {
		panic(err)
	}
}
