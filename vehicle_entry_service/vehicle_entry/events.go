package main

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type VehicleEntryEvent struct {
	ID            string `json:"id"`
	VehiclePlate  string `json:"vehicle_plate"`
	EntryDateTime string `json:"entry_date_time"`
}

func publishEntryEvent(ch *amqp.Channel, ctx context.Context, event VehicleEntryEvent) error {
	body, err := json.Marshal(event)

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		"vehicle_entry",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}
