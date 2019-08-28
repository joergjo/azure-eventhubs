package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
)

// EventHubReceiver is used to receive messages from an Event Hub
type EventHubReceiver struct {
	hub     *eventhub.Hub
	handler func(c context.Context, e *eventhub.Event) error
}

var receiveLog = struct {
	sync.Mutex
	items map[string]int
}{
	items: make(map[string]int),
}

// NewEventHubReceiver creates a new EventHubReceiver instance
func NewEventHubReceiver(connectionString string, handler func(context.Context, *eventhub.Event) error) (*EventHubReceiver, error) {
	hub, err := eventhub.NewHubFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	receiver := &EventHubReceiver{hub: hub, handler: handler}
	return receiver, nil
}

// Listen starts accepting messages from an Event Hub
func (r *EventHubReceiver) Listen() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	runtimeInfo, err := r.hub.GetRuntimeInformation(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, partitionID := range runtimeInfo.PartitionIDs {
		fmt.Printf("Receiving on partition %s...\n", partitionID)
		_, err := r.hub.Receive(ctx, partitionID, createHandler(partitionID, r.handler), eventhub.ReceiveWithConsumerGroup("$Default"), eventhub.ReceiveFromTimestamp(now))
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops an EventHubReceiver
func (r *EventHubReceiver) Stop() error {
	err := r.hub.Close(context.Background())
	fmt.Println("Receiver stopped.")
	fmt.Println("***")
	fmt.Println("Dumping statistics...")
	for k, v := range receiveLog.items {
		fmt.Printf("Partition %s received %d messages.\n", k, v)
	}
	return err
}

func logReceive(key string) {
	receiveLog.Lock()
	receiveLog.items[key] = receiveLog.items[key] + 1
	receiveLog.Unlock()
}

func createHandler(partitionID string, innerHandler func(context.Context, *eventhub.Event) error) func(context.Context, *eventhub.Event) error {
	return func(c context.Context, e *eventhub.Event) error {
		logReceive(partitionID)
		return innerHandler(c, e)
	}
}
