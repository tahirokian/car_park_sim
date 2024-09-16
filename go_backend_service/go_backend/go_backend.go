package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func getRabbitmqUrl() string {
	rabbitmqAddress := os.Getenv("RABIITMQ_ADDR")
	rabbitmqPort := os.Getenv("RABIITMQ_PORT")
	rabbitmqUrl := fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitmqAddress, rabbitmqPort)
	return rabbitmqUrl
}

func getRedisUrl() string {
	redisAddress := os.Getenv("REDIS_ADDR")
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl := fmt.Sprintf("%s:%s", redisAddress, redisPort)
	return redisUrl
}

func main() {
	start := time.Now()

	var conn *amqp.Connection
	var err error
	var rdb *redis.Client

	for {
		conn, err = amqp.Dial(getRabbitmqUrl())
		if err != nil {
			log.Println("Failed to connect to RabbitMQ, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to RabbitMQ successfully!")
		break
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	defer ch.Close()

	rdb = redis.NewClient(&redis.Options{
		Addr:     getRedisUrl(),
		Password: "",
		DB:       0,
	})

	defer rdb.Close()

	log.Println("Go backend service setup took: ", time.Since(start))

	// Consume vehicle entry events
	go func() {
		msgs, _ := ch.Consume("vehicle_entry", "", true, false, false, false, nil)
		for msg := range msgs {
			var event VehicleEntryEvent
			json.Unmarshal(msg.Body, &event)
			recordEntryEvent(rdb, event)
		}
	}()

	// Consume vehicle exit events
	go func() {
		msgs, _ := ch.Consume("vehicle_exit", "", true, false, false, false, nil)
		for msg := range msgs {
			var event VehicleExitEvent
			json.Unmarshal(msg.Body, &event)
			recordExitEventAndSummary(rdb, event)
		}
	}()

	// Keep service running
	select {}
}
