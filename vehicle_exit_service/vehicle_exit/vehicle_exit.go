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
	"github.com/redis/go-redis/v9"
)

func randomVehiclePlate(rdb *redis.Client) string {
	plates, err := rdb.SMembers(context.Background(), "entry_plates").Result()

	// No vehicle has been registered.
	if err != nil || len(plates) == 0 {
		return ""
	}

	// 80% probability to pick from an existing plate
	if rand.Float64() < 0.8 {
		return plates[rand.Intn(len(plates))]
	} else {
		// 20% probability to generate a random, unmatched plate
		// For the exercise, XYZ is used as prefix. This should be randomly generated.
		return fmt.Sprintf("XYZ-%d", rand.Intn(10))
	}
}

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

type metrics struct {
	numberOfProcessedEvents prometheus.Counter
	startupDelay            prometheus.Gauge
	eventProcessingDelay    prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		numberOfProcessedEvents: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "vehicle_exit",
				Name:      "total_processed_exit_events",
				Help:      "Total number of processed vehicle exit events.",
			},
		),
		startupDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "vehicle_exit",
				Name:      "startup_delay",
				Help:      "Startup delay for vehicle exit service.",
			},
		),
		eventProcessingDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "vehicle_exit",
				Name:      "exit_event_processing_delay",
				Help:      "Entry event processing delay for vehicle exit service.",
			},
		),
	}

	reg.MustRegister(m.numberOfProcessedEvents)
	reg.MustRegister(m.startupDelay)
	reg.MustRegister(m.eventProcessingDelay)
	return m
}

func setupPromethusEndpoint(reg *prometheus.Registry) {
	vehicleExitAddress := os.Getenv("VEHICLE_EXIT_ADDR")
	vehicleExitPort := os.Getenv("VEHICLE_EXIT_PORT")
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	http.Handle("/metrics", promHandler)
	log.Printf("Prometheus metrics available at http://%s:%s/metrics\n",
		vehicleExitAddress, vehicleExitPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", vehicleExitPort), nil))
}

func processExitEvents(ch *amqp.Channel, ctx context.Context, rdb *redis.Client, m *metrics) {
	for {
		eventProcessingStart := time.Now()

		vehicle_plate := randomVehiclePlate(rdb)

		if vehicle_plate != "" {

			exitEvent := VehicleExitEvent{
				ID:           fmt.Sprintf("event-%d", rand.Intn(10000)),
				VehiclePlate: vehicle_plate,
				ExitDateTime: time.Now().UTC().Format(time.RFC3339),
			}

			for {
				err := publishExitEvent(ch, ctx, exitEvent)

				if err != nil {
					log.Println("Failed to publish exit event, retrying in 5 seconds...")
					time.Sleep(5 * time.Second)
					continue
				}

				log.Printf("Published exit event: %+v\n", exitEvent)
				break
			}
		}
		eventProcessingDurationSec := time.Since(eventProcessingStart).Seconds()
		log.Println("Exit event processing took: ", eventProcessingDurationSec)

		m.numberOfProcessedEvents.Inc()
		m.eventProcessingDelay.Set(float64(eventProcessingDurationSec))

		// Add a delay for next exit event
		time.Sleep(2 * time.Second)
	}
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

	// Connect to Redis to fetch entered vehicle plates
	rdb := redis.NewClient(&redis.Options{
		Addr:     getRedisUrl(),
		Password: "",
		DB:       0,
	})

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	m := NewMetrics(reg)
	go setupPromethusEndpoint(reg)

	setupDuration := time.Since(start)
	m.startupDelay.Set(setupDuration.Seconds())

	log.Println("Vehicle exit service setup took: ", time.Since(start))

	go processExitEvents(ch, ctx, rdb, m)

	select {}
}
