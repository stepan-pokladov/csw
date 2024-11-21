package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/stepan-pokladov/csw/client/internal/client"
)

// Endpoints for the visit and activity APIs
const (
	visitEndpoint    = "/api/visit/v1"
	activityEndpoint = "/api/activity/v1"
)

func main() {
	cc := os.Getenv("CLIENT_COUNT")
	if cc == "" {
		cc = "1"
	}
	clientCount, err := strconv.Atoi(cc)
	if err != nil {
		log.Fatalf("Invalid CLIENT_COUNT: %v", err)
	}
	ri := os.Getenv("REPORT_INTERVAL")
	if ri == "" {
		ri = "5"
	}
	reportInterval, err := time.ParseDuration(ri + "s")
	if err != nil {
		log.Fatalf("Invalid REPORT_INTERVAL: %v", err)
	}

	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		cc = "http://localhost:8080"
	}

	var wg sync.WaitGroup
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		c := client.NewClient(serverUrl+visitEndpoint, serverUrl+activityEndpoint, reportInterval)
		go c.Run(&wg)
	}
	wg.Wait()
}
