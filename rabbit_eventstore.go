package onewallet

import (
	"errors"

	framework "github.com/djansyle/go-eventsource"
	rabbit "github.com/djansyle/go-rabbit"
)

// RabbitEventStoreClient is the client for rabbitmq communication for the eventstore
type RabbitEventStoreClient struct {
	client rabbit.Sender
}


// NewRabbitEventStoreClient Creates a new instance of the rabbit client with the
func NewRabbitEventStoreClient(opt *rabbit.CreateClientOption) (framework.EventStore, error) {
	client, err := rabbit.CreateClient(opt, requestFormatter)
	if err != nil {
		return nil, err
	}

	return &RabbitEventStoreClient{client: client}, nil
}

// RetrieveEvents gets a list of events from the event store
func (c *RabbitEventStoreClient) RetrieveEvents(opts *framework.RetrieveEventsOption) (events []framework.Event, err error) {
	if c.client == nil {
		return nil, errors.New("rabbitmq not initialized")
	}


	err = c.client.Send(Request{
		Action: "Events",
		Data:   *opts,
	}, &events)

	return events, err
}

// CreateEvent send a new event to be saved to the event store
func (c *RabbitEventStoreClient) CreateEvent(event framework.Event) (newEvent framework.Event, err error) {
	err = c.client.Send(Request{
		Action: "CreateEvent",
		Data:   event,
	}, &newEvent)

	return newEvent, err
}

// CreateSnapshot sends a message to the eventstore to create a snapshot
func (c *RabbitEventStoreClient) CreateSnapshot(opts *framework.CreateSnapshotOption) (success bool, err error) {
	err = c.client.Send(Request{
		Action: "CreateSnapshot",
		Data:   *opts,
	}, &success)

	return success, err
}

// RetrieveSnapshot retrieves a snapshot from the eventstore
func (c *RabbitEventStoreClient) RetrieveSnapshot(opts *framework.RetrieveSnapshotOption) (snapshot framework.Snapshot, err error) {
	err = c.client.Send(Request{
		Action: "Snapshot",
		Data:   *opts,
	}, &snapshot)

	return snapshot, err
}
