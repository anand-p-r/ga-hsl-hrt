package main

import (
	"os"
	//"time"
	//"fmt"
	//"strings"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

var routeInfo []routeData

var configRoutes []string
var configSigns map[string]string
var file *os.File

func enableLogging() {
	//timeStr := time.Now().Format("2006-01-02T15:04:05")
	//fmt.Printf("time-%v", timeStr)
	//timeStr = strings.ReplaceAll(timeStr, "T", "-")
	//timeStr = strings.ReplaceAll(timeStr, ":", "-")
	//logFile := LOGFILE + "-" + timeStr + ".log"
	logFile := LOGFILE + ".log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)
	log.Info("Logging to a file-", logFile)
	log.SetLevel(log.TraceLevel)
}

func getConfig() {
	viper.SetConfigName("config-file")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Panic("Fatal error config file: ", err)
	}

	// Get routes
	configRoutes = viper.GetStringSlice(ROUTES)
	configSigns = viper.GetStringMapString(SIGNS)

	log.Info("routes - ", configRoutes)
	log.Info("dest - ", configSigns)

	return
}

func main() {

	enableLogging()
	getConfig()

	routeInf := buildRouteData(JUPPERI1)
	routeInfo = append(routeInfo, routeInf)

	routeInf = buildRouteData(JUPPERI2)
	routeInfo = append(routeInfo, routeInf)

	routeInf = buildRouteData(JUPPERI3)
	routeInfo = append(routeInfo, routeInf)

	listenAndServe()
}
