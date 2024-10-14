package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type VehicleExitEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	ExitDateTime string `json:"exit_date_time"`
}

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

func publishExitEvent(ch *amqp.Channel, ctx context.Context, event VehicleExitEvent) error {
	body, err := json.Marshal(event)

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		"vehicle_exit",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}
