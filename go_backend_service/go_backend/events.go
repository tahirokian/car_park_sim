package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type VehicleEntryEvent struct {
	ID            string `json:"id"`
	VehiclePlate  string `json:"vehicle_plate"`
	EntryDateTime string `json:"entry_date_time"`
}

type VehicleExitEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	ExitDateTime string `json:"exit_date_time"`
}

type VehicleParkingSummary struct {
	VehiclePlate  string `json:"vehicle_plate"`
	EntryDateTime string `json:"entry_date_time"`
	ExitDateTime  string `json:"exit_date_time"`
	Duration      string `json:"duration"`
}

var ctx = context.Background()

func recordEntryEvent(rdb *redis.Client, event VehicleEntryEvent) {
	start := time.Now()

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Failed to jsonify entry event %+v", event)
	}

	rdb.HSet(ctx, event.VehiclePlate, "entry_time", event.EntryDateTime)
	rdb.HSet(ctx, event.VehiclePlate, "parking_enter_data", data)
	rdb.SAdd(ctx, "entry_plates", event.VehiclePlate)
	log.Printf("Recorded vehicle entry event to redis: %v\n", event)

	log.Println("Saving entry event to redis took: ", time.Since(start))
}

func recordExitEventAndSummary(rdb *redis.Client, event VehicleExitEvent) {
	start := time.Now()

	entryTime, err := rdb.HGet(ctx, event.VehiclePlate, "entry_time").Result()

	if err == redis.Nil {
		log.Printf("No entry record found for vehicle %v. Ignoring...\n", event.VehiclePlate)
	} else {
		data, json_err := json.Marshal(event)
		if json_err != nil {
			log.Fatalf("Failed to jsonify exit event %+v", event)
		}

		rdb.HSet(ctx, event.VehiclePlate, "parking_exit_data", data)
		log.Printf("Recorded vehicle exit event to redis: %v\n", event)

		exit_time, _ := time.Parse(time.RFC3339, event.ExitDateTime)
		entry_time, _ := time.Parse(time.RFC3339, entryTime)
		duration := exit_time.Sub(entry_time)

		// Call REST API to store summary
		sendToPythonServer(event.VehiclePlate, entryTime, event.ExitDateTime, duration)

		log.Println("Saving exit entry to redis and python server took: ", time.Since(start))
	}
}
