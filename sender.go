package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

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
	sender := &EventHubSender{hub: hub}
	return sender, nil
}

// Send sends a message
func (s *EventHubSender) send(num int, key *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	message := fmt.Sprintf("Message %d with partition key %s", num, *key)
	log.Printf("Sending '%s'\n", message)
	e := eventhub.NewEventFromString(message)
	e.PartitionKey = key
	err := s.hub.Send(ctx, e)
	return err
}

// Send sends multiple messages message
func (s *EventHubSender) Send(limit int) error {
	for i := 0; i < limit; i++ {
		uuid, _ := uuid.NewRandom()
		key := uuid.String()
		if err := s.send(i, &key); err != nil {
			return err
		}
	}
	return nil
}
