package main

import (
	"context"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
)

// EventHubReceiver is used to receive messages from an Event Hub
type EventHubReceiver struct {
	hub     *eventhub.Hub
	handler func(context.Context, *eventhub.Event) error
}

// NewEventHubReceiver creates a new EventHubReceiver instance
func NewEventHubReceiver(connectionString string, handler func(context.Context, *eventhub.Event) error) (*EventHubReceiver, error) {
	hub, err := eventhub.NewHubFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	receiver := EventHubReceiver{hub: hub, handler: handler}
	return &receiver, nil
}

// Listen starts accepting messages from an Event Hub
func (r *EventHubReceiver) Listen() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	runtimeInfo, err := r.hub.GetRuntimeInformation(ctx)
	if err != nil {
		return err
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		_, err := r.hub.Receive(ctx, partitionID, r.handler, eventhub.ReceiveWithConsumerGroup("$Default"))
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops an EventHubReceiver
func (r *EventHubReceiver) Stop() error {
	return r.hub.Close(context.Background())
}
