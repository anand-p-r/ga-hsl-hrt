package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

/* Helper function - For REST responses */
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}

func listenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/getRoute/{route}", GetRouteHandler)

	srv := &http.Server{
		Addr: "localhost:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	time.Sleep(time.Duration(3600) * time.Second)
}

func GetRouteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("Route: %v\n", vars["route"])

	var restResp []restResponse
	var respR restResponse

	if vars["route"] == ROUTE215 {
		fmt.Println("\n1\n")
		for _, val := range routeInfo {
			fmt.Println("\n2\n")
			for _, val2 := range val.arrDepDetails {
				fmt.Printf("\n3 - %v\n", val2.headSign)
				if (strings.Contains(val2.headSign, HEADSIGN215_LEP)) ||
					(strings.Contains(val2.headSign, HEADSIGN215_LAH)) {
					fmt.Println("\n4\n")
					respR.StopCode = val.stopDetails.code
					respR.HeadSign = val2.headSign
					respR.ScheduledDep = timeFromSeconds(val2.scheduledDeparture)
					respR.RealArr = timeFromSeconds(val2.realtimeArrival)

					restResp = append(restResp, respR)
				}
			}
		}
	}
	fmt.Printf("REST RESP - %v", restResp)
	respondWithJSON(w, http.StatusOK, restResp)
	return
}

func timeFromSeconds(seconds float64) (calcTime time.Time) {

	currTime := time.Now()

	loc, _ := time.LoadLocation("Local")

	calcTime = time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 0, 0, int(seconds), 0, loc)

	return

}
