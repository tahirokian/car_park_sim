package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type VehicleExitEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	ExitDateTime string `json:"exit_date_time"`
}

func publishExitEvent(ch *amqp.Channel, ctx context.Context, event VehicleExitEvent) error {
	body, err := json.Marshal(event)

	if err != nil {
		log.Fatalf("Failed to jsonify the event: %v", err)
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
