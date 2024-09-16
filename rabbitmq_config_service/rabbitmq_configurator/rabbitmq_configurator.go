package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

	// Declare the queues for vehicle entry and exit events
	_, err = ch.QueueDeclare("vehicle_entry", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare vehicle_entry queue: %v", err)
	}

	_, err = ch.QueueDeclare("vehicle_exit", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare vehicle_exit queue: %v", err)
	}

	log.Println("Successfully declared vehicle_entry and vehicle_exit queues...")
	log.Println("Total execution time: ", time.Since(start))
}
