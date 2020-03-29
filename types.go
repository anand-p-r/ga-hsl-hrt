package main

import (
	"time"
)

const url string = "https://api.digitransit.fi/routing/v1/routers/hsl/index/graphql"

// Sharing the server certs from hassio. These are autorenewed by Lets Encrypt Addon. How about that!!
const (
	SERVERCERT string = "/usr/share/hassio/ssl/fullchain.pem"
	SERVERKEY string = "/usr/share/hassio/ssl/privkey.pem"
	LOGFILE string = "./ga-hsl-hrt"
)

// Intent matching strings
const (
  DESTONLY string = "Destination-Only"
  BUSDEST string = "Bus-Destination"
)

// Configuration keys
const (
	ROUTES string = "routes"
	SIGNS  string = "callSignToHeadsign"
)

const (
  
  // 215 only towards Leppa
  JUPPERI1        string = "HSL:2143202"
  
  // Towards Helsinki
  JUPPERI3        string = "HSL:2143218"
  
  // 215 towards Lähderanta, 214 towards Leppävaara, 548 towards Tapiola, 321 towards Vanhakartano, 565 towards Espoontori
  JUPPERI2        string = "HSL:2143217"
)


const (

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

/*

{
  "fulfillmentText": "This is a text response",
  "fulfillmentMessages": [
    {
      "card": {
        "title": "card title",
        "subtitle": "card text",
        "imageUri": "https://example.com/images/example.png",
        "buttons": [
          {
            "text": "button text",
            "postback": "https://example.com/path/for/end-user/to/follow"
          }
        ]
      }
    }
  ],
  "source": "example.com",
  "payload": {
    "google": {
      "expectUserResponse": true,
      "richResponse": {
        "items": [
          {
            "simpleResponse": {
              "textToSpeech": "this is a simple response"
            }
          }
        ]
      }
    },
    "facebook": {
      "text": "Hello, Facebook!"
    },
    "slack": {
      "text": "This is a text response for Slack."
    }
  },
  "outputContexts": [
    {
      "name": "projects/project-id/agent/sessions/session-id/contexts/context-name",
      "lifespanCount": 5,
      "parameters": {
        "param-name": "param-value"
      }
    }
  ],
  "followupEventInput": {
    "name": "event name",
    "languageCode": "en-US",
    "parameters": {
      "param-name": "param-value"
    }
  }
}

*/

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