package onewallet

import (
	"testing"
	"time"

	"github.com/djansyle/go-rabbit"
)

type Command string
type Query string
type Root string

func (c *Command) Create(data interface{}) (string, *rabbit.ApplicationError) {
	return "someid", nil
}

func (q *Query) Accounts(data interface{}) ([]map[string]string, *rabbit.ApplicationError) {
	return []map[string]string{{"username": "john"}, {"username": "johndoe"}}, nil
}

func (r *Root) CreateNothing(input interface{}) (bool, *rabbit.ApplicationError) {
	return true, nil
}

func TestService(t *testing.T) {
	service, err := NewService(NewServiceOpt{URL: "amqp://guest:guest@localhost:5672/", Name: "Testing", Version: 1, Root: new(Root), Query: new(Query), Command: new(Command)})
	if err != nil {
		t.Fatalf("expecting error to be `nil` but got %q", err)
	}

	go service.Start()

	client, err := rabbit.CreateClient(&rabbit.CreateClientOption{URL: "amqp://guest:guest@localhost:5672/", Queue: "Testing", TimeoutRequest: 2 * time.Second})
	if err != nil {
		t.Fatalf("expect `err` to be `nil` but got %q", err)
	}

	var createResult string
	client.Send(Request{
		Action: "Command.Create",
		Data:   map[string]string{},
	}, &createResult)

	if createResult != "someid" {
		t.Fatalf("expecting value to be %q but got %q", "someid", createResult)
	}

	var queryResult []map[string]string
	client.Send(Request{
		Action: "Query.Accounts",
		Data:   map[string]string{},
	}, &queryResult)

	if queryResult[0]["username"] != "john" {
		t.Fatalf("expecting value to be to contain a property name %q to have a value of %q but got %v", "username", "john", queryResult)
	}

	rootResult := false
	client.Send(Request{
		Action: "CreateNothing",
		Data:   map[string]string{},
	}, &rootResult)

	if rootResult != true {
		t.Fatalf("expecting value to be `true` but got %v", rootResult)
	}

}
