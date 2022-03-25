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

	// set microservice listener threads
	src.SetAlphaListenerThread()
}