package main

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var config *Config
var l *TestLogger
var Locate2UAccessToken string

func invokeSync() {

	// syncCustomers()
	for range time.Tick(time.Second * 3600 * 24 * time.Duration(config.Environment.SyncCycle)) {
		syncCustomers()
	}
}

func invokeGetAccessToken() {

	// Get access token per hour(actually 50 min)
	for range time.Tick(time.Second * 60 * 50) {
		Locate2UAccessToken = getAccessToken(config.Locate2U)
	}
}

func main() {

	config = GetConfig()
	Locate2UAccessToken = getAccessToken(config.Locate2U)
	go invokeGetAccessToken()
	go invokeSync()

	// Logger for testing purposes only - to be replaced with production logger
	l = NewTestLogger()

	switch config.Environment.LogLevel {
	case 1:
		log.SetLevel(log.WarnLevel)
	case 2:
		log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

	log.SetFormatter(&log.JSONFormatter{})

	if config.Environment.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	l.Debug("Starting...")

	// Routes...
	router := gin.Default()

	// Test sync with this endpoint
	router.POST("/sync", doSync)
	router.POST("/trip/:fulfillmentid", trip_to)
	router.Run("localhost:8000")

	// log.Fatal(r.Run(fmt.Sprintf("%s:%d", config.Environment.Host, config.Environment.Port)))
}
