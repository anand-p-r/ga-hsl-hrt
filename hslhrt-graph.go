/*
hslhrt-graph.go

Handles the GraphQL interface functionality towards HSL API.
- For every bus stop configured in the configuration file, retrieve all buses including their 
  scheduled and realtime departure times and store them in route information structure.
- Route information structures are later used by the webser handler to respond with bus timings.
*/

package main

import (
	"context"
	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
	"strings"
	"sort"
)

// Universal GraphQL client variable
var graphClient *graphql.Client = graphql.NewClient(url)

// Sort structure and functions for scheduled arrival time
type aDSlice []routeArrDepDetails
func (aD aDSlice) Len() int { 
	return len(aD) 
}

func (aD aDSlice) Less(i, j int) bool { 
	return aD[i].scheduledDeparture < aD[j].scheduledDeparture 
}

func (aD aDSlice) Swap(i, j int) {
	aD[i], aD[j] = aD[j], aD[i]
}

/*
getRoutesFromStop: Retrieves routes from HSL API over the GraphQL interface
Include the used GraphQL query that can also be verified at this link:
https://api.digitransit.fi/graphiql/hsl. 
Input: gtfsId that uniquely identifies a bus stop
Output: Returns the route data structure with arrival and departure times of routes
from the bus stop.
*/
func getRoutesFromStop(gtfsId string) (routeInfo routeData) {

	req := graphql.NewRequest(`query ($id: String!) {
		stop (id: $id) {
			name
			code
			routes {
			  shortName
			  patterns{
				headsign
			  }
			}
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

	var stopDet stopStruct
	var arrivalDeparture []routeArrDepDetails
	var routeSigns []routeHeadSigns

	stopDet.gtfsId = gtfsId

	if stop := respMap["stop"]; stop != nil {
		stopDetails := stop.(map[string]interface{})

		for key, val := range stopDetails {
			switch key {
			case "name":
				stopDet.name = val.(string)
			case "lat":
				stopDet.latitude = val.(float64)
			case "lon":
				stopDet.longitude = val.(float64)
			case "code":
				stopDet.code = val.(string)
			case "stoptimesWithoutPatterns":
				stopTimes := stopDetails["stoptimesWithoutPatterns"].([]interface{})
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
			default:
				// Do nothing
			}
		}

		// Sort routes based on scheduled arrival time
		sort.Sort(aDSlice(arrivalDeparture))
		routeInfo.stopDetails = stopDet
		routeInfo.arrDepDetails = arrivalDeparture

		// Now map the routes to the headsigns
		// Headsign in the stop pattern is the destination of the route/bus
		// Headsign in the routes struct contains both start and dest.
		// So compare the headsign of the stop pattern with headsigns of the routes to find 
		// the route for that particulartime in the stop pattern
		// NOTE: Headsign Names in the routes are shorter than head sign names in the stop pattern

		// Build the route signs structure
		var routeSign routeHeadSigns
		if routeNames := stopDetails["routes"]; routeNames != nil {
			routeNames := routeNames.([]interface{})

			for _, val := range routeNames {
				patMap := val.(map[string]interface{})
				if patterns := patMap["patterns"]; patterns != nil {
					patterns := patMap["patterns"].([]interface{})

					var signs []string
					for _, val := range patterns {
						hsMap := val.(map[string]interface{})
						signs = append(signs, hsMap["headsign"].(string))
					}
					routeSign.headsigns = signs
				}
				routeSign.routeName = patMap["shortName"].(string)
				routeSigns = append(routeSigns, routeSign)
			}
			log.Debug("routeSigns-", routeSigns)
		}

		// Compare routes and arrivalDep details to match destination headsigns
		for indx, arrDep := range routeInfo.arrDepDetails {
			found := false
			for _, rt := range routeSigns {
				for _, hs := range rt.headsigns {
					if strings.Contains(strings.ToLower(arrDep.headSign), strings.ToLower(hs)) {
						routeInfo.arrDepDetails[indx].route = rt.routeName
						found = true
						break
					}
				}

				if found {
					break
				}
			}
			log.Debug("arrDep - ", arrDep)
		}
	}

	return
}

/*
buildRouteData: Function that builds route information for a stop configured in 
configuration file.
Input: gtfsId that uniquely identifies a bus stop
Output: Returns the route data structure with arrival and departure times of routes
from the bus stop.
*/
func buildRouteData(gtfsId string) (routeInfo routeData) {

	routeInfo = getRoutesFromStop(gtfsId)

	log.Info("Route Info - ", routeInfo)

	return
}
