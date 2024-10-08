package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

type metrics struct {
	numberOfProcessedEvents prometheus.Counter
	startupDelay            prometheus.Gauge
	eventProcessingDelay    prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		numberOfProcessedEvents: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "vehicle_entry",
				Name:      "total_processed_events",
				Help:      "Total number of processed vehicle entry events.",
			},
		),
		startupDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "vehicle_entry",
				Name:      "startup_delay",
				Help:      "Startup delay for vehicle entry service.",
			},
		),
		eventProcessingDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "vehicle_entry",
				Name:      "event_processing_delay",
				Help:      "Event processing delay for vehicle entry service.",
			},
		),
	}

	reg.MustRegister(m.numberOfProcessedEvents)
	reg.MustRegister(m.startupDelay)
	reg.MustRegister(m.eventProcessingDelay)
	return m
}

func setupPromethusEndpoint(reg *prometheus.Registry) {
	vehicleEntryAddress := os.Getenv("VEHICLE_ENTRY_ADDR")
	vehicleEntryPort := os.Getenv("VEHICLE_ENTRY_PORT")
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	http.Handle("/metrics", promHandler)
	log.Printf("Prometheus metrics available at http://%s:%s/metrics\n",
		vehicleEntryAddress, vehicleEntryPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", vehicleEntryPort), nil))
}

func processEnrtyEvents(ch *amqp.Channel, ctx context.Context, m *metrics) {
	for {
		eventProcessingStart := time.Now()

		entryEvent := VehicleEntryEvent{
			ID:            fmt.Sprintf("event-%d", rand.Intn(10000)),
			VehiclePlate:  randomVehiclePlate(),
			EntryDateTime: time.Now().UTC().Format(time.RFC3339),
		}

		for {
			err := publishEntryEvent(ch, ctx, entryEvent)

			if err != nil {
				log.Println("Failed to publish enter event, retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("Published entry event: %+v\n", entryEvent)
			break
		}
		eventProcessingDurationSec := time.Since(eventProcessingStart).Seconds()
		log.Printf("Entry event processing took: %v\n", eventProcessingDurationSec)

		m.numberOfProcessedEvents.Inc()
		m.eventProcessingDelay.Set(float64(eventProcessingDurationSec))

		// Add a delay before next event.
		time.Sleep(2 * time.Second)
	}
}

func main() {
	startTime := time.Now()

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

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	m := NewMetrics(reg)
	go setupPromethusEndpoint(reg)

	_, declareErr := ch.QueueDeclare("vehicle_entry", true, false, false, false, nil)
	if declareErr != nil {
		log.Fatalf("Failed to declare vehicle_entry queue: %v", declareErr)
	}

	setupDuration := time.Since(startTime)
	m.startupDelay.Set(setupDuration.Seconds())
	log.Println("Vehicle entry service setup took: ", setupDuration)

	go processEnrtyEvents(ch, ctx, m)

	select {}
}
