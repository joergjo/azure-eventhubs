# azure-eventhubs

An expanded version of the [Azure Event Hubs Go SDK's basic sample](https://github.com/Azure/azure-event-hubs-go). 

## Building
```
cd $GOPATH/src/github.vom/joergjo/azure-eventhubs
go get go get -u github.com/Azure/azure-event-hubs-go
go build
```

## Running
The tool can be either run in sender or receiver mode

- `-connection-string` EventHub connection string including a SAS for either listening or sending
- `-receiver` Run tool as receiver, otherwise it acts as sender (`false`by default)
- `-send-count` Number of messages to send (ignored in receiver mode)

The tool runs until all messages have been sent in sender mode, or until it has been signalled to quit (e.g. by pressing CTRL+C) in receiver mode.
