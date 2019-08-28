package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	eventhub "github.com/Azure/azure-event-hubs-go"
)

var receiver bool
var connectionString string
var sendCount int

func init() {
	flag.BoolVar(&receiver, "receiver", false, "Receiver mode")
	flag.StringVar(&connectionString, "connection-string", "", "Event Hub sender connection string")
	flag.IntVar(&sendCount, "send-count", 100, "Number of messages to send")
}

func main() {
	flag.Parse()
	if receiver {
		receive()
	} else {
		send()
	}
}

func send() {
	ehSender, err := NewEventHubSender(connectionString)
	if err != nil {
		log.Fatalf("Failed to create EventHub client for sending: %s\n", err)
	}
	fmt.Printf("Sending %d messages...\n", sendCount)
	if err = ehSender.Send(sendCount); err != nil {
		log.Fatalf("Failed to send all messages: %s\n", err)
	}
	fmt.Println("Done sending.")
}

func receive() {
	ehReceiver, err := NewEventHubReceiver(connectionString, eventHandler)
	if err != nil {
		log.Fatalf("Failed to create EventHub client for receiving: %s\n", err)
	}
	if err = ehReceiver.Listen(); err != nil {
		log.Fatalf("Failed to start receiver: %s\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	if err = ehReceiver.Stop(); err != nil {
		log.Fatalf("Failed to stop receiver: %s\n", err)
	}
}

func eventHandler(c context.Context, event *eventhub.Event) error {
	fmt.Println(string(event.Data))
	return nil
}
