/*
ga-hsl-hrt.go

Main file for the webserver implementation
- Configuration parsing and updating is done here.
- Logging is enabled.
- Main entry point for the webserver
*/

package main

import (
	"os"
	"sync"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

// Main structure that holds routes retrieved from HSL API
var routeInfo []routeData

// Configuration parameters
var configRoutes []string
var configSigns map[string]string
var configStopGtfsIds []string
var listeningPort string
var clientCaCert string
var serverCert string
var serverKey string

// Logfile
var file *os.File

// Waitgroup for go routines
var wg sync.WaitGroup

/* 
enableLogging: Function to enable logging to a file. Utilises logrus package
to define different logging levels.   
*/
func enableLogging(logFile string) {
	logFile = logFile + ".log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)
	log.Info("Logging to a file-", logFile)
	log.SetLevel(log.TraceLevel)
}

/*
getConfig: Function to parse and extract configuration parameters from
config-file.json. All parameters are obligatory. Missing parameters will
cause a non recoverable panic!
*/
func getConfig() {
	viper.SetConfigName("config-file")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Panic("Fatal error config file: ", err)
	}

	// Get configuration parameters
	configRoutes = viper.GetStringSlice(ROUTES)
	configSigns = viper.GetStringMapString(SIGNS)
	listeningPort = viper.GetString(PORT)
	logFile := viper.GetString(LOGFILE)
	clientCaCert = viper.GetString(CLIENTCERT)
	serverCert = viper.GetString(SERVERCERT)
	serverKey = viper.GetString(SERVERKEY)
	configStopGtfsIds = viper.GetStringSlice(STOPGTFSIDS)

	if logFile == "" {
		panic("No logfile defined!")
	} else {
		enableLogging(logFile)
	}

	if len(configRoutes) == 0 {
		log.Panic("No routes defined!")
	}

	if len(configSigns) == 0 {
		log.Panic("No headsigns defined!")
	}

	if len(configStopGtfsIds) == 0 {
		log.Panic("No stopGtfsIds defined!")
	}

	if listeningPort == "" {
		log.Panic("Server port not defined in config file!")
	}

	if serverCert == "" {
		log.Panic("Server certificate location not defined in config file!")
	}

	if serverKey == "" {
		log.Panic("Server Key location not defined in config file!")
	}

	if clientCaCert == "" {
		log.Panic("Client certificate location not defined in config file!")
	}

	log.Info("routes - ", configRoutes)
	log.Info("callsigns - ", configSigns)
	log.Info("stopgtfsids - ", configStopGtfsIds)
	log.Info("port - ", listeningPort)
	log.Info("serverCert - ", serverCert)
	log.Info("serverKey - ", serverKey)
	log.Info("clientCert - ", clientCaCert)

	return
}

/*
main: Exatly what it says!! Main function
*/
func main() {

	// Read configuration information
	getConfig()

	// Populate internal structures from Graphql response
	for _, stop := range(configStopGtfsIds) {
		routeInf := buildRouteData(stop)
		routeInfo = append(routeInfo, routeInf)
	}

	// Start the webserver
	listenAndServe()

	// Do not exit as long as the webser is running
	log.Debug("Waiting for pending go routines to end...")
	wg.Wait()
	log.Debug("Everything is done. Ending Main!!")
	file.Close()
}

