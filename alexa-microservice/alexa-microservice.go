package main

import (
	"alexa-microservice/config"
	"alexa-microservice/src"
	"fmt"
)

// Runs the micro-service
func main() {

	// quit if config.json contains errors
	if err := config.GetErrStatus(); err != nil {
		fmt.Println(err)
		return
	}

	//// create go routines for all microservice listener threads

	// WolframAlpha Queries
	go src.SetAlphaListenerThread()
	
	// Twxt to Speech
	go src.SetTTSListenerThread()
	
	// Speech to Text
	go src.SetSTTListenerThread()
	
	// Alexa (Speech to Speech)
	go src.SetAlexaListenerThread()

	select {} // listen until program termination
}