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

var senderConnectionString string
var sendCount int
var receiverConnectionString string

func init() {
	flag.StringVar(&senderConnectionString, "sndConnStr", "", "Event Hub sender connection string")
	flag.IntVar(&sendCount, "sndCnt", 100, "Number of messages to send")
	flag.StringVar(&receiverConnectionString, "rcvConnStr", "", "Event Hub receiver connection string")
}

func main() {
	flag.Parse()
	if senderConnectionString != "" {
		send()
	}

	if receiverConnectionString != "" {
		receive()
	}
}

func send() {
	ehSender, err := NewEventHubSender(senderConnectionString)
	if err != nil {
		log.Fatalf("Failed to create EventHub client for sending: %s\n", err)
	}
	fmt.Printf("Sending %d messages...\n", sendCount)
	if err = ehSender.Send(sendCount); err != nil {
		log.Fatalf("Failed to send all messages: %s", err)
	}
	fmt.Println("Done sending.")
}

func receive() {
	ehReceiver, err := NewEventHubReceiver(receiverConnectionString, eventHandler)
	if err != nil {
		log.Fatalf("Failed to create EventHub client for receiving: %s\n", err)
	}
	if err = ehReceiver.Listen(); err != nil {
		log.Fatalf("Failed to start listeners: %s\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	if err = ehReceiver.Stop(); err != nil {
		log.Fatalf("Failed to stop listeners: %s\n", err)
	}
	fmt.Println("Receiver has shut down.")
}

func eventHandler(c context.Context, event *eventhub.Event) error {
	fmt.Println(string(event.Data))
	return nil
}
