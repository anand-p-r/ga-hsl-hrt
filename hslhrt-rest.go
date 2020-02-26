package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
						fmt.Printf("key - %v\n", key3)
						fmt.Printf("val - %v\n\n", val3)
					}
				}
			}
		}
	}

	fmt.Printf("\ntarg-%v", target)
	fmt.Printf("\nERR-%v", err)

	resp.Body.Close()

	jsonBody = []byte(`{"query": "{route(id: \"HSL:1009\") {shortName longName}}"}`)
	resp2, err := http.Post(url, header, bytes.NewBuffer(jsonBody))

	err = json.NewDecoder(resp2.Body).Decode(&target)

	for key, val := range target {
		fmt.Printf("key - %v\n", key)
		fmt.Printf("val - %v\n\n", val)
	}

	fmt.Printf("\ntarg-%v", target)
	fmt.Printf("\nERR-%v", err)

	resp2.Body.Close()
}
