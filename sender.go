package main

import (
	"context"
	"fmt"
	"log"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
)

// EventHubSender is used to send messages to an Event Hub
type EventHubSender struct {
	hub *eventhub.Hub
}

// NewEventHubSender creates a new EventHubSender
func NewEventHubSender(connectionString string) (*EventHubSender, error) {
	hub, err := eventhub.NewHubFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	sender := EventHubSender{hub: hub}
	return &sender, nil
}

// Send sends a message
func (s *EventHubSender) send(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	log.Printf("Sending '%s'", message)
	err := s.hub.Send(ctx, eventhub.NewEventFromString(message))
	return err
}

// Send sends multiple messages message
func (s *EventHubSender) Send(limit int) error {
	for i := 0; i < limit; i++ {
		message := fmt.Sprintf("Message %d", i)
		if err := s.send(message); err != nil {
			return err
		}
	}
	return nil
}
