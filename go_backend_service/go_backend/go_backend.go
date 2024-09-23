package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	var wg sync.WaitGroup
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

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	m := NewMetrics(reg)
	go setupPromethusEndpoint(reg)

	setupDuration := time.Since(start)
	m.startupDelay.Set(setupDuration.Seconds())

	log.Println("Go backend service setup took: ", setupDuration)

	wg.Add(1)

	// Consume vehicle entry events
	go func() {
		defer wg.Done()
		msgs, _ := ch.Consume("vehicle_entry", "", true, false, false, false, nil)

		for msg := range msgs {
			entryEventStart := time.Now()

			var event VehicleEntryEvent
			json.Unmarshal(msg.Body, &event)
			recordEntryEvent(rdb, event)

			entryEventProcessingDuration := time.Since(entryEventStart)
			m.eventProcessingDelay.WithLabelValues("entry_event").
				Set(entryEventProcessingDuration.Seconds())
			m.eventProcessingDelayHist.WithLabelValues("entry_event").
				Observe(entryEventProcessingDuration.Seconds())
			m.numberOfProcessedEvents.WithLabelValues("entry_event").Inc()
		}
	}()

	wg.Add(1)

	// Consume vehicle exit events
	go func() {
		defer wg.Done()

		msgs, _ := ch.Consume("vehicle_exit", "", true, false, false, false, nil)

		for msg := range msgs {
			exitEventStart := time.Now()
			var event VehicleExitEvent
			json.Unmarshal(msg.Body, &event)
			recordExitEventAndSummary(rdb, event)

			exitEventProcessingDuration := time.Since(exitEventStart)
			m.eventProcessingDelay.WithLabelValues("exit_event").
				Set(exitEventProcessingDuration.Seconds())
			m.eventProcessingDelayHist.WithLabelValues("exit_event").
				Observe(exitEventProcessingDuration.Seconds())
			m.numberOfProcessedEvents.WithLabelValues("exit_event").Inc()
		}
	}()

	wg.Wait()
}
