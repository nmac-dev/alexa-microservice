package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Struct representation of config.json
type Config struct {
	
	AppName string `json:"AppName"`

	// WolframAlpha (Alpha)
	WolframAlpha struct {
		ApiURI		string `json:"apiURI"`
		AppPath		string `json:"appPath"`
		AppID		string `json:"appID"`
	} `json:"WolframAlpha"`

	// Microsoft Cognitive Services (MCS)
	MicrosoftCS struct {
		Region 		string `json:"region"`
		TTS			string `json:"tts"`
		STT			string `json:"stt"`
		STTLang 	string `json:"sttLang"`
		SrvPath 	string `json:"srvPath"`
		AzureKey	string `json:"azureKey"`
	} `json:"MicrosoftCS"`

	// Extensible Markup Language (XML)
	Xml struct {
		Version 	string `xml:"version"`
		Encoding	string `xml:"encoding"`
		Lang 		string `xml:"lang"`
	} `xml:"Xml"`
	// Speech Synthesis Markup Language (SSML)
	Ssml struct {
		Voice struct {
			Name 	string `xml:"name"`
		} `xml:"voice"`
	} `xml:"Ssml"`
}

// private package struct, ensures only one config exists
var cfgAlexa Config
var cfgErr error = nil

// Loads the config.json file & parses it to a private config struct (loads once)
func getConfig() (*Config) {

	// ensure json config is loaded only one time
	if cfgAlexa.AppName == "" {
		const filename string = "config/config.json"

		file, err := os.Open(filename)
		if err != nil {
			cfgErr = fmt.Errorf("could not load: %s \nE: %s", filename, err)
		}

		// parse json -> config struct
		err = json.NewDecoder(file).Decode(&cfgAlexa)
		
		if err != nil {
			cfgErr = fmt.Errorf("failed to parse JSON from: %s \nE:%s", filename, err)
		}
	}
	return &cfgAlexa // never nil, initialised in package scope
}

////	Getters		////

func GetAppName()		string { return getConfig().AppName }
func GetErrStatus()		error  { getConfig(); return cfgErr }

func GetAlphaApiURI()	string { return getConfig().WolframAlpha.ApiURI	 }
func GetAlphaAppPath()	string { return getConfig().WolframAlpha.AppPath }
func GetAlphaAppID()	string { return getConfig().WolframAlpha.AppID	 }

func GetMcsRegion()		string { return getConfig().MicrosoftCS.Region	 }
func GetMcsTTS()		string { return getConfig().MicrosoftCS.TTS		 }
func GetMcsSTT()		string { return getConfig().MicrosoftCS.STT		 }
func GetMcsSTTLang()	string { return getConfig().MicrosoftCS.STTLang	 }
func GetMcsSrvPath()	string { return getConfig().MicrosoftCS.SrvPath	 }
func GetMcsAzureKey()	string { return getConfig().MicrosoftCS.AzureKey }

func GetXmlVersion()	string { return getConfig().Xml.Version		}
func GetXmlEncoding()	string { return getConfig().Xml.Encoding	}
func GetXmlLang()		string { return getConfig().Xml.Lang		}
func GetSsmlVoiceName()	string { return getConfig().Ssml.Voice.Name	}

// Return the config in a JSON representation
func ToString() string {

	if GetAppName() == "" || GetErrStatus() != nil{
		return "{}: config is empty"	
	}
	cfgString, _ := json.MarshalIndent(cfgAlexa, "\n", "\t")

	return string(cfgString)
}
