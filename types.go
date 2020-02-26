package main

import (
	"time"
)

const url string = "https://api.digitransit.fi/routing/v1/routers/hsl/index/graphql"

const (
	ROUTE215 string = "215"
	ROUTE214 string = "214"
	ROUTE548 string = "548"
	ROUTE565 string = "565"
	ROUTE321 string = "321"
)

// 215 only towards Leppa
const (
	JUPPERI1        string = "HSL:2143202"
	HEADSIGN215_LEP string = "Leppävaara"
)

// 215 towards Lähderanta, 214 towards Leppävaara, 548 towards Tapiola, 321 towards Vanhakartano, 565 towards Espoontori
const (
	JUPPERI2        string = "HSL:2143217"
	HEADSIGN214     string = "Leppävaara"
	HEADSIGN215_LAH string = "Lähderanta"
	HEADSIGN548     string = "Tapiola"
	HEADSIGN565     string = "Espoontori"
	HEADSIGN321     string = "Vanhakartano"
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
