package main

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"log"
)

// Universal client variable
var graphClient *graphql.Client = graphql.NewClient(url)

func getRoutesFromStop(gtfsId string) (arrivalDeparture []routeArrDepDetails) {

	req := graphql.NewRequest(`query ($id: String!) {
		stop (id: $id) {
			name
		    stoptimesWithoutPatterns {
				scheduledArrival
		  		realtimeArrival
		  		arrivalDelay
		  		scheduledDeparture
		  		realtimeDeparture
		  		departureDelay
		  		realtime
		  		realtimeState
		  		serviceDay
				headsign
			}
		}
	}`)

	req.Var("id", gtfsId)
	ctx := context.Background()

	var respMap map[string]interface{}

	if err := graphClient.Run(ctx, req, &respMap); err != nil {
		log.Fatal(err)
	}

	if stop := respMap["stop"]; stop != nil {
		stopDetails := stop.(map[string]interface{})

		if stopTimes := stopDetails["stoptimesWithoutPatterns"]; stopTimes != nil {

			stopTimes := stopTimes.([]interface{})

			for _, stopAD := range stopTimes {
				var arrDep routeArrDepDetails

				arDepTimes := stopAD.(map[string]interface{})

				// Scheduled arrival time
				if item := arDepTimes["scheduledArrival"]; item != nil {
					arrDep.scheduledArrival = item.(float64)
				}

				// Realtime arrival time
				if item := arDepTimes["realtimeArrival"]; item != nil {
					arrDep.realtimeArrival = item.(float64)
				}

				// Arrival delay time
				if item := arDepTimes["arrivalDelay"]; item != nil {
					arrDep.arrivalDelay = item.(float64)
				}

				// Scheduled departure time
				if item := arDepTimes["scheduledDeparture"]; item != nil {
					arrDep.scheduledDeparture = item.(float64)
				}

				// Real time departure time
				if item := arDepTimes["realtimeDeparture"]; item != nil {
					arrDep.realtimeDeparture = item.(float64)
				}

				// Departure delay time
				if item := arDepTimes["departureDelay"]; item != nil {
					arrDep.departureDelay = item.(float64)
				}

				// Realtime data available
				if item := arDepTimes["realtime"]; item != nil {
					arrDep.realtime = item.(bool)
				}

				// Realtime state available
				if item := arDepTimes["realtimeState"]; item != nil {
					arrDep.realtimeState = item.(string)
				}

				// Headsign
				if item := arDepTimes["headsign"]; item != nil {
					arrDep.headSign = item.(string)
				}

				arrivalDeparture = append(arrivalDeparture, arrDep)
			}
		}
	}

	return
}

func getStop(gtfsId string) (stopS stopStruct) {

	req := graphql.NewRequest(`query ($id: [String]!) {
		stops(ids: $id){
			gtfsId
			name
			code 
			lat 
			lon
			}
		}`)

	req.Var("id", gtfsId)

	ctx := context.Background()

	var resp map[string]interface{}

	if err := graphClient.Run(ctx, req, &resp); err != nil {
		log.Fatal(err)
	}

	mainKey := "stops"
	if ok := resp[mainKey]; ok == nil {
		fmt.Println("Key-%v not found", mainKey)
		return
	}

	respStops := resp[mainKey].([]interface{})

	for _, slice := range respStops {
		sliceVal := slice.(map[string]interface{})
		for key, val := range sliceVal {
			switch key {
			case "name":
				stopS.name = val.(string)
			case "lat":
				stopS.latitude = val.(float64)
			case "lon":
				stopS.longitude = val.(float64)
			case "code":
				stopS.code = val.(string)
			case "gtfsId":
				stopS.gtfsId = val.(string)
			default:
				// Do nothing
			}
		}
	}
	return
}

func buildRouteData(gtfsId string) (routeInfo routeData) {

	routeInfo.stopDetails = getStop(gtfsId)
	routeInfo.arrDepDetails = getRoutesFromStop(gtfsId)

	fmt.Printf("\nRoute Info - %v\n", routeInfo)

	return
}
