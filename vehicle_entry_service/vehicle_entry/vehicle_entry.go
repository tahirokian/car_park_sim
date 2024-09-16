package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func randomVehiclePlate() string {
	/*
		Vehicle plates should be generated the following way.
			letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			plate := fmt.Sprintf("%c%c%c-%d%d%d", letters[rand.Intn(26)], letters[rand.Intn(26)],
			letters[rand.Intn(26)], rand.Intn(10), rand.Intn(10), rand.Intn(10))
		For the demo, I am using a very small subset so that the functionality can be verified.
	*/
	plate := fmt.Sprintf("ABC-%d", rand.Intn(10))
	return plate
}

func getRabbitmqUrl() string {
	rabbitmqAddress := os.Getenv("RABIITMQ_ADDR")
	rabbitmqPort := os.Getenv("RABIITMQ_PORT")
	rabbitmqUrl := fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitmqAddress, rabbitmqPort)
	return rabbitmqUrl
}

func main() {
	start := time.Now()

	var conn *amqp.Connection
	var err error

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Vehicle entry service setup took: ", time.Since(start))

	for {
		eventProcessingStart := time.Now()

		entryEvent := VehicleEntryEvent{
			ID:            fmt.Sprintf("event-%d", rand.Intn(10000)),
			VehiclePlate:  randomVehiclePlate(),
			EntryDateTime: time.Now().UTC().Format(time.RFC3339),
		}

		for {
			err = publishEntryEvent(ch, ctx, entryEvent)

			if err != nil {
				log.Println("Failed to publish enter event, retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("Published entry event: %+v\n", entryEvent)
			break
		}
		log.Println("Entry event processing took: ", time.Since(eventProcessingStart))

		// Add a delay before next event.
		time.Sleep(2 * time.Second)
	}
}
