/*
http-server.go

HTTPS Webserver Handling with mutual TLS
- Parse and extract required parameters from Webhook Request from Google Assistant.
- Format response into Webhook Response for Google Assistant.
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
)

/* 
respondWithError: Helper function - For REST responses with error status
*/
func respondWithError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}

/* 
respondWithJSON: Helper function - For REST responses with valid JSON payload
*/
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}

/*
timeFromSeconds: Helper function - To convert time from float64 to locale based 
formatted string.
*/
func timeFromSeconds(seconds float64) (calcTime time.Time) {

	currTime := time.Now()

	loc, _ := time.LoadLocation("Local")

	calcTime = time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 0, 0, int(seconds), 0, loc)

	return

}

/*
listenAndServe: Gathers the server and client certificates before starting the 
webserver. Server is started in a go routine so its non-blocking.
TODO: Can extend this further to listen to signals over channels to 
shutdown and restart webserver - for e.g. when runtime configuration updates are done
This also means runtime configuration updates have to be supported first :)
*/
func listenAndServe() {

	r := mux.NewRouter()
	r.HandleFunc("/getRoute", GetRouteHandler).Methods("POST")

	// Read google client certificates for mTLS
	// curl https://pki.goog/gsr2/GTS1O1.crt | openssl x509 -inform der >> google-certs\ca-cert.pem
	// curl https://pki.goog/gsr2/GSR2.crt | openssl x509 -inform der >> google-certs\ca-cert.pem
	// Also generate client side certificate for the host from where curl will be issued for testing and use that cert and key in curl command
	// curl -X GET <https:domain:port/getRoute/215> --cert ./localhost.pem --key ./localhost.out -v

	caCerts, err := ioutil.ReadFile(clientCaCert)
	if err != nil {
		log.Error(err)
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCerts)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()	

	addr := ":" + listeningPort
	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		// IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
		TLSConfig: tlsConfig, // Client side certificates
	}

	// Run our server in a goroutine so that it doesn't block.
	wg.Add(1)
	go func() {
		if err := srv.ListenAndServeTLS(serverCert, serverKey); err != nil {
			log.Error(err)
		}
		wg.Done()
	}()
}

/*
extractPostParams: Extracts JSON body parameters from Webhook Request received from
Dialogflow.
*/
func extractPostParams (r *http.Request) (request string, route string, destination string) {

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)

	if err != nil {
		log.Error("Error decoding json body in request - ", err)
		return
	}

	// We are trying to extract this part of json from the POST body
	// "queryResult":{"queryText":"route 215",
	// 		"parameters":{"route":[54,8],"place-attraction":"tapiola"},......"intent":{"name":"...","displayName":"Bus-Destination"},
	//		"parameters":{"place-attraction":"sello"},......"intent":{"name":"...","displayName":"Destination-Only"},

	_, ok := t["queryResult"]
	if ok {
		qResult := t["queryResult"].(map[string]interface{})

		// Check intent type and extract parameters
		_, ok:= qResult["intent"]

		if ok {
			intent := qResult["intent"].(map[string]interface{})

			_, ok := intent["displayName"]

			if ok {
				rcvdIntent := intent["displayName"].(string)

				switch strings.ToLower(rcvdIntent) {
				case strings.ToLower(DESTONLY):
					request = DESTONLY
				case strings.ToLower(BUSDEST):
					request = BUSDEST
				default:
					log.Error("Unsupported intent received-", rcvdIntent)
				} 
			}
		}

		_, ok = qResult["parameters"]

		if ok {
			parameters := qResult["parameters"].(map[string]interface{})

			// Extract route parameter which is available only if the intent is bus-dest
			if request == BUSDEST {
				_, ok := parameters["route"]

				// route is a list which can be [2,1,4] or [21,4] or [2,14] or [214]
				// Best option is to convert each number to string and concatenate them
	
				if ok {
					routeList := parameters["route"].([]interface{})
					for i:=0;i<len(routeList);i++ {
						st := fmt.Sprintf("%.0f", routeList[i].(float64))
						route = route + st
					}
	
					// Intelligent??: Take only first 3 numbers to ignore user conversational errors
					if len(route) > 3 {
						route = route[:3]
					}	
				} else {
					log.Debug("parameters[route] - ", parameters["route"])
					return
				}
			}

			// Extract destination
			_, ok := parameters["place-attraction"]
			if ok {
				destination = parameters["place-attraction"].(string)
			} else {
				log.Debug("parameters[place-attraction] - ", parameters["place-attraction"])
			}
				
		} else {
			log.Debug("qResult[parameters] - ", qResult["parameters"])
			return
		}
	} else {
		log.Debug("t[queryResult] - ", t["queryResult"])
		return
	}

	return
}

/*
GetBusDestinationHandler: Handler to extract bus and destination details from route structurees,
based on given bus and destination.
Formats them into a string slice with scheduled/realtime departure timings
*/
func GetBusDestinationHandler(route string, headSign string) (routes []string){
	var routeString string
	for _, rtInfo := range routeInfo {
		found := false
		routeString = "Leaves from " + 
			rtInfo.stopDetails.name + 
			" (" + rtInfo.stopDetails.code + ")" + 
			" at "
		for _, arrDep := range rtInfo.arrDepDetails {
			if strings.Contains(arrDep.route, route) {
				if strings.Contains(strings.ToLower(arrDep.headSign), strings.ToLower(headSign)) {
					// Bingo!
					if found {
						routeString = routeString + ", "
					}

					found = true
					if arrDep.realtime {
						routeString = routeString + timeFromSeconds(arrDep.realtimeDeparture).Format("15:04")
					} else {
						routeString = routeString + timeFromSeconds(arrDep.scheduledDeparture).Format("15:04")
					}
				}
			}
		}

		if found {
			routes = append(routes, routeString)
		}
	}

	return
}

/*
GetDestinationHandler: Handler to extract bus and destination details from route structures,
based on given destination.
Formats them into a string slice with scheduled/realtime departure timings
*/
func GetDestinationHandler(headSign string) (routes []string) {
	// Populate the GA Webhook Response Struct
	var buses []string
	var times []float64
	for _, rtInfo := range routeInfo {
		for _, arrDep := range rtInfo.arrDepDetails {
			if strings.Contains(strings.ToLower(arrDep.headSign), strings.ToLower(headSign)) {
				found := false
				for indx, bus := range buses {
					if strings.Contains(arrDep.route, bus) {
						found = true
						// We already have found this bus before for this dest, but now we have a new time. Append it
						if arrDep.realtime {
							routes[indx] = routes[indx] + "," + timeFromSeconds(arrDep.realtimeDeparture).Format("15:04")
						} else {
							routes[indx] = routes[indx] + "," + timeFromSeconds(arrDep.scheduledDeparture).Format("15:04")
						}	
						break
					}	
				}

				if !found {
					// New bus found for the given destination. Create an entry in routes.
					buses = append(buses, arrDep.route)
					routeString := "Bus " + arrDep.route + 
						" leaves from " + 
						rtInfo.stopDetails.name + 
						" (" + rtInfo.stopDetails.code + ")" + 
						" at "
					var fTime float64
					if arrDep.realtime {
						routeString = routeString + timeFromSeconds(arrDep.realtimeDeparture).Format("15:04")
						fTime = arrDep.realtimeDeparture
					} else {
						routeString = routeString + timeFromSeconds(arrDep.scheduledDeparture).Format("15:04")
						fTime = arrDep.scheduledDeparture
					}

					times = append(times, fTime) 
					routes = append(routes, routeString)	
				}
			}
		}
	}

	// Sort the routes based on times
	log.Debug("Times - ", times)
	log.Debug("Routes - ", routes)
	for i:=0;i<len(times)-1;i++ {
		for j:=0;j<(len(times)-i-1);j++ {
			if times[j] > times[i] {
				// swap routes
				routes[i], routes[j] = routes[j], routes[i]
				times[i], times[j] = times[j], times[i]
			}
		}
	}
	log.Debug("Aft Times - ", times)
	log.Debug("Aft Routes - ", routes)

	return
}

/*
GetRouteHandler: Webser main handler that invokes other functions based on available 
bus+destination or destination.
Formats the string slice into Webhook Response format for Dialogflow.
*/
func GetRouteHandler(w http.ResponseWriter, r *http.Request) {

	request, route, destination := extractPostParams(r)
	log.Info("WebHook Req for request - ", request, " route - ", route, " to destination - ", destination)

	var gaWebHkResp gaWebHookResponse
	var items []itemStruct	

	headSign, ok := configSigns[strings.ToLower(destination)]

	if !ok {
		log.Error("Destination could not be mapped - ", destination)
		gaWebHkResp.FulfillmentText = "Destination could not be mapped"
	
		item := itemStruct{
			SimpleResponse: simpleRespStruct{
				TextToSpeech: "Sorry, but destination could not be mapped! Please retry.",
			},
		}

		items = append(items, item)

		richResp := richResponseStruct{
			Items: items,
		}
	
		gaWebHkResp.Payload = payloadStruct{
			Google: googleStruct{
				ExpectUserResponse: true,
				RichResponse: richResp,
			},
		}
	
		log.Error("WebHook ERROR RESP - ", gaWebHkResp)
		respondWithJSON(w, http.StatusOK, gaWebHkResp)
		return
	}

	var routes []string

	switch request {
	case BUSDEST: 
		routes = GetBusDestinationHandler(route, headSign)
	case DESTONLY:
		routes = GetDestinationHandler(headSign)
	default:
		log.Error("Unsupported handler type - ", request, "Internal error!!")
	}

	if len(routes) == 0 {
		log.Error("Route-", route, " Destination-", destination, " mismatch")
		gaWebHkResp.FulfillmentText = "No routes to provided destination"
	
		item := itemStruct{
			SimpleResponse: simpleRespStruct{
				TextToSpeech: "Sorry, but no routes were found! Please retry.",
			},
		}

		items = append(items, item)

		richResp := richResponseStruct{
			Items: items,
		}
	
		gaWebHkResp.Payload = payloadStruct{
			Google: googleStruct{
				ExpectUserResponse: true,
				RichResponse: richResp,
			},
		}
	
		log.Error("WebHook ERROR RESP - ", gaWebHkResp)
		respondWithJSON(w, http.StatusOK, gaWebHkResp)
		return
	}

	// Only two simple responses are expected
	if len(routes) > 2 {
		routes = routes[:2]
	}
	
	for _, rt := range routes {
		item := itemStruct{
			SimpleResponse: simpleRespStruct{
				TextToSpeech: rt,
			},
		}
	
		items = append(items, item)
	}

	richResp := richResponseStruct{
		Items: items,
	}
	
	gaWebHkResp.FulfillmentText = "Here are upcoming buses for route " + route
	
	gaWebHkResp.Payload = payloadStruct{
		Google: googleStruct{
			ExpectUserResponse: true,
			RichResponse: richResp,
		},
	}

	log.Info("WebHook RESP - ", gaWebHkResp)
	respondWithJSON(w, http.StatusOK, gaWebHkResp)

	return

}