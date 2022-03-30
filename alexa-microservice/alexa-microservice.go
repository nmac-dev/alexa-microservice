package main

import (
	"alexa-microservice/config"
	"alexa-microservice/services"
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
	go services.SetAlphaListenerThread()
	
	// Twxt to Speech
	go services.SetTTSListenerThread()
	
	// Speech to Text
	go services.SetSTTListenerThread()
	
	// Alexa (Speech to Speech)
	go services.SetAlexaListenerThread()

	select {} // listen until program termination
}