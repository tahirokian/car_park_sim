package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getWebBackendUrl() string {
	webBackendAddress := os.Getenv("WEBBACKEND_ADDR")
	webBackendPort := os.Getenv("WEBBACKEND_PORT")
	webBackendUrl := fmt.Sprintf("http://%s:%s/upload", webBackendAddress, webBackendPort)
	return webBackendUrl
}

func sendToPythonServer(plate, entry, exit string, duration time.Duration) {
	start := time.Now()

	summaryEvent := VehicleParkingSummary{
		VehiclePlate:  plate,
		EntryDateTime: entry,
		ExitDateTime:  exit,
		Duration:      duration.String(),
	}

	jsonData, _ := json.Marshal(summaryEvent)

	resp, err := http.Post(getWebBackendUrl(), "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("Error sending data to Python server: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Summary saved to file for plate %v\n", plate)
	} else {
		log.Printf("Summary saving failed with http status code: %v\n", resp.StatusCode)
	}

	log.Println("Sending summary to python server took: ", time.Since(start))
}
