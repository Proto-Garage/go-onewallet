package onewallet

import (
	"testing"
	"time"

	framework "github.com/djansyle/go-eventsource"
	rabbit "github.com/djansyle/go-rabbit"
)

func assertNilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expecting error to be nil but got %v", err)
	}
}

type incomingRequest struct {
	Request string
	Data    interface{}
}

type EventStoreMock rabbit.Service

func (*EventStoreMock) CreateSnapshot(i framework.CreateSnapshotOption) (interface{}, *rabbit.ApplicationError) {
	return true, nil
}

func (*EventStoreMock) CreateEvent(i framework.Event) (framework.Event, *rabbit.ApplicationError) {
	i.ID = rabbit.RandomID()
	return i, nil
}

func (*EventStoreMock) Snapshot(i framework.RetrieveSnapshotOption) (framework.Snapshot, *rabbit.ApplicationError) {
	var snapshot = framework.Snapshot{AggregateID: i.AggregateID, AggregateType: i.AggregateType, AggregateVersion: 1, State: []byte("1")}
	return snapshot, nil
}

func (*EventStoreMock) Events(i framework.RetrieveEventsOption) ([]framework.Event, *rabbit.ApplicationError) {
	return []framework.Event{framework.Event{AggregateID: i.AggregateID, Type: "MockedEvent"}}, nil
}

func startEventstore(t *testing.T) {
	server, err := rabbit.CreateServer("amqp://guest:guest@localhost:5672/", "EventStore")
	assertNilError(t, err)

	server.RegisterName("", new(EventStoreMock))

	go server.Serve()
}

// TestEventStore will test the IO of the rabbitmq
func TestEventStore(t *testing.T) {
	startEventstore(t)
	rabbitEventStore, err := NewRabbitEventStoreClient(&rabbit.CreateClientOption{URL: "amqp://guest:guest@localhost:5672/", Queue: "EventStore", TimeoutRequest: time.Second})
	assertNilError(t, err)

	id := rabbit.RandomID()
	result, err := rabbitEventStore.RetrieveEvents(&framework.RetrieveEventsOption{AggregateID: id})

	assertNilError(t, err)

	if len(result) == 0 {
		t.Fatal("expecting return `result` to be more than 1")
	}

	if result[0].AggregateID != id {
		t.Fatalf("expecting arrived event an event[0] to have an aggregate id of %s instead got %s", id, result[0].AggregateID)
	}

	event, err := rabbitEventStore.CreateEvent(framework.Event{AggregateID: rabbit.RandomID(), AggregateType: 2})
	if event.ID == "" {
		t.Fatal("expecting `ID` field to exists instead got `nil`")
	}

	success, err := rabbitEventStore.CreateSnapshot(&framework.CreateSnapshotOption{AggregateID: rabbit.RandomID(), AggregateVersion: 1, AggregateType: 100})
	if success != true {
		t.Fatalf("expecting `success` to be `true` but got %v", success)
	}

	snapshot, err := rabbitEventStore.RetrieveSnapshot(&framework.RetrieveSnapshotOption{AggregateType: 100, AggregateID: id})
	if snapshot.AggregateID != id {
		t.Fatalf("expecting aggregate id to equal to %q instead got %q", id, snapshot.AggregateID)
	}
}
