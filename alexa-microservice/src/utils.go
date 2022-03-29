package src

import (
	"alexa-microservice/config"
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
)

//// Hosts utility functions, variables, & structs used in data representation

const (
	resPath		= "res/"
	resPerms	= fs.FileMode(0666); // -rw-rw-rw- | read & write 
)

var (
	xmlHeader = fmt.Sprintf(
		`<?xml version="%s" encoding="%s"?>%c`, 
		config.GetXmlVersion(), config.GetXmlEncoding(), '\n',
	)
	// request header for MCS azure key	
	hAzureKey = rHeader{Key: "Ocp-Apim-Subscription-Key", Value: config.GetMcsAzureKey() }
)

type (
	
	// JSON text
	JsonText struct {
		Data string `json:"text"`
	}

	// JSON speech
	JsonSpeech struct {
		Data []byte `json:"speech"`
	}

	// request header key/value
	rHeader struct {
		Key 	string
		Value 	string
	}
	
	// Speech Synthesis Markup Language (SSML) struct to be sent to microsoft cognitive services
	SSMLSpeak struct {
		XMLName		xml.Name	`xml:"speak"`
		Version		string		`xml:"version,attr"`
		Lang		string		`xml:"xml:lang,attr"`
		Voice		SSMLVoice	`xml:"voice"`
	} // <speak> <voice> <speak/>

	SSMLVoice struct {
		XMLName		xml.Name	`xml:"voice"`
		Lang		string		`xml:"xml:lang,attr"`
		Name		string		`xml:"name,attr"`
		CharData	[]byte		`xml:",chardata"`
	} // <voice>

	// MCS json response fields: {"RecognitionStatus","DisplayText","Offset","Duration"}
	// DisplayText is the field required for JsonText
	MCSJsonRsp struct {
		RecognitionStatus 	string
		DisplayText			string
		Offset				int
		Duration			int
	}
)

// Helper function to ignore double valued returns
func singular(val []byte, _ interface{}) *([]byte) { return &val }

// Outputs the data to the specific filename within the resource directory
func writeToResFile(filename string, data []byte) {

	// skip if resPath exists
	if _, err := os.Stat(resPath); err != nil{ 
		// create resource directory
		err = os.Mkdir(resPath, resPerms)
		if err != nil {
			err = fmt.Errorf(
				"cannot create '%s' directory with perms '%d': update app permissions \nE: '%s'",
				resPath, resPerms, err,
			)
			fmt.Println(err)
			panic(err)
		}
	}
	
	// output data to file
	err := os.WriteFile(resPath + filename, data, resPerms)
	if err != nil {
		err = fmt.Errorf(
			"failed to create '%s' in '%s' with perms '%d': check app permissions \nE: %s", 
			filename, resPath, resPerms, err,
		)
		fmt.Println(err)
		panic(err)
	}
}