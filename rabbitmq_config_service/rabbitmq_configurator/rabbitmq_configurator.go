package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func getRabbitmqUrl() string {
	rabbitmqAddress := os.Getenv("RABIITMQ_ADDR")
	rabbitmqPort := os.Getenv("RABIITMQ_PORT")
	rabbitmqUrl := fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitmqAddress, rabbitmqPort)
	return rabbitmqUrl
}

type metrics struct {
	startupDelay prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		startupDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "queue_configurator",
				Name:      "startup_delay",
				Help:      "Startup delay for queue configurator service.",
			},
		),
	}

	reg.MustRegister(m.startupDelay)
	return m
}

func setupPromethusEndpoint(reg *prometheus.Registry) {
	queueConfigAddress := os.Getenv("QUEUE_CONFIG_ADDR")
	queueConfigPort := os.Getenv("QUEUE_CONFIG_PORT")
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	http.Handle("/metrics", promHandler)
	log.Printf("Prometheus metrics available at http://%s:%s/metrics\n",
		queueConfigAddress, queueConfigPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", queueConfigPort), nil))
}

func setupQueues(ch *amqp.Channel) {
	// Declare the queues for vehicle entry and exit events
	_, err := ch.QueueDeclare("vehicle_entry", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare vehicle_entry queue: %v", err)
	}

	_, err = ch.QueueDeclare("vehicle_exit", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare vehicle_exit queue: %v", err)
	}

	log.Println("Successfully declared vehicle_entry and vehicle_exit queues...")
}

func main() {
	start := time.Now()

	var conn *amqp.Connection
	var err error

	conn, err = amqp.Dial(getRabbitmqUrl())
	if err != nil {
		log.Fatalf("Failed to open rabbitmq connection: %v", err)
	}

	log.Println("Connected to RabbitMQ successfully!")

	defer conn.Close()

	ch, chErr := conn.Channel()
	if chErr != nil {
		log.Fatalf("Failed to open a channel: %v", chErr)
	}
	defer ch.Close()

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	m := NewMetrics(reg)
	go setupPromethusEndpoint(reg)

	setupQueues(ch)

	setupDuration := time.Since(start)
	m.startupDelay.Set(setupDuration.Seconds())
	log.Println("Time to setup rabbitmq queues: ", setupDuration)

	select {}
}
