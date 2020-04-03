# azure-eventhubs-sender-receiver

An expanded version of the [Azure Event Hubs Go SDK's basic sample](https://github.com/Azure/azure-event-hubs-go). 

## Building (using Go 1.14 or higher)
```
cd $GOPATH/src/github.vom/joergjo/azure-eventhubs
go build
```

## Running
The application can be either run in sender or receiver mode

- `-connection-string` EventHub connection string including a SAS for either listening or sending
- `-receiver` Run tool as receiver, otherwise it acts as sender (`false`by default)
- `-send-count` Number of messages to send (ignored in receiver mode)

The application runs until all 
- messages have been sent (sender mode)
- it has been signalled to quit, e.g. by pressing CTRL+C (receiver mode)
