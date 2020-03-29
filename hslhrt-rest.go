package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func restClient() {

	//Url
	header := "application/json"

	// Create a json body with request params
	jsonBody := []byte(`{"query": "{stop(id: \"HSL:1040129\") {name lat lon wheelchairBoarding}}"}`)
	resp, err := http.Post(url, header, bytes.NewBuffer(jsonBody))

	var target map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&target)

	for key, val := range target {
		if key == "data" {
			val1 := val.(map[string]interface{})
			for key2, val2 := range val1 {
				if key2 == "stop" {
					val2 := val2.(map[string]interface{})
					for key3, val3 := range val2 {
						log.Debug("key - ", key3)
						log.Debug("val - ", val3)
					}
				}
			}
		}
	}

	log.Debug("targ - ", target)
	log.Debug("ERR - ", err)

	resp.Body.Close()

	jsonBody = []byte(`{"query": "{route(id: \"HSL:1009\") {shortName longName}}"}`)
	resp2, err := http.Post(url, header, bytes.NewBuffer(jsonBody))

	err = json.NewDecoder(resp2.Body).Decode(&target)

	for key, val := range target {
		log.Debug("key - ", key)
		log.Debug("val - ", val)
	}

	log.Debug("targ - ", target)
	log.Debug("ERR - ", err)

	resp2.Body.Close()
}
