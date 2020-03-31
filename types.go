/*
types.go

- Structure types and constants used across the application.

*/

package main

const url string = "https://api.digitransit.fi/routing/v1/routers/hsl/index/graphql"

// Intent matching strings
const (
  DESTONLY string = "Destination-Only"
  BUSDEST string = "Bus-Destination"
)

// Configuration keys
const (
	ROUTES      string = "routes"
  SIGNS       string = "callSignToHeadsign"
  PORT        string = "port"
  LOGFILE     string = "logFile"
  CLIENTCERT  string = "clientCert"
  SERVERCERT  string = "serverCert"
  SERVERKEY   string = "serverKey"
  STOPGTFSIDS string = "stopGtfsIds"
)

// A bus's arrival/departure details.
type routeArrDepDetails struct {
	scheduledArrival   float64
	realtimeArrival    float64
	arrivalDelay       float64
	scheduledDeparture float64
	realtimeDeparture  float64
	departureDelay     float64
	realtime           bool
	realtimeState      string
	headSign           string
	route              string
}

// A bus's headsigns source/destination
type routeHeadSigns struct {
	routeName string
	headsigns []string
}

// Main structure that holds all stops and buses from the stop.
type routeData struct {
	stopDetails   stopStruct
	arrDepDetails []routeArrDepDetails
}

// Structure that holds the bus stop details
type stopStruct struct {
	gtfsId    string
	name      string
	code      string
	latitude  float64
	longitude float64
}

// Structure for Webhook Response to Dialogflow.
type simpleRespStruct struct {
	TextToSpeech string `json:"textToSpeech"`
}

// Structure for Webhook Response to Dialogflow.
type itemStruct struct {
	SimpleResponse simpleRespStruct `json:"simpleResponse"`
}

// Structure for Webhook Response to Dialogflow.
type suggestionStruct struct {
	Title string `json:"title"`
}

// Structure for Webhook Response to Dialogflow.
type richResponseStruct struct {
	Items       []itemStruct       `json:"items"`
	Suggestions []suggestionStruct `json:"suggestions"`
}

// Structure for Webhook Response to Dialogflow.
type googleStruct struct {
	ExpectUserResponse bool               `json:"expectUserResponse"`
	RichResponse       richResponseStruct `json:"richResponse"`
}

// Structure for Webhook Response to Dialogflow.
type payloadStruct struct {
	Google             googleStruct       `json:"google"`
}

// Structure for Webhook Response to Dialogflow.
type gaWebHookResponse struct {
	FulfillmentText  string        `json:"fulfillmentText"`
	Payload          payloadStruct `json:"payload"`
}