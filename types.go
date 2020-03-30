package main

import (
	"time"
)

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

type routeHeadSigns struct {
	routeName string
	headsigns []string
}

type routeData struct {
	stopDetails   stopStruct
	arrDepDetails []routeArrDepDetails
}

type stopStruct struct {
	gtfsId    string
	name      string
	code      string
	latitude  float64
	longitude float64
}

type restResponse struct {
	StopCode     string    `json: StopCode`
	HeadSign     string    `json: HeadSign`
	ScheduledDep time.Time `json: ScheduledDep`
	RealArr      time.Time `json: RealArr`
}

type simpleRespStruct struct {
	TextToSpeech string `json:"textToSpeech"`
}

type itemStruct struct {
	SimpleResponse simpleRespStruct `json:"simpleResponse"`
}

type suggestionStruct struct {
	Title string `json:"title"`
}

// Webhooks are limited to two simple responses.
// Let's also limit ourselves to one suggestion.
type richResponseStruct struct {
	Items       []itemStruct       `json:"items"`
	Suggestions []suggestionStruct `json:"suggestions"`
}

type googleStruct struct {
	ExpectUserResponse bool               `json:"expectUserResponse"`
	RichResponse       richResponseStruct `json:"richResponse"`
}

type payloadStruct struct {
	Google             googleStruct       `json:"google"`
}

type gaWebHookResponse struct {
	FulfillmentText  string        `json:"fulfillmentText"`
	Payload          payloadStruct `json:"payload"`
}