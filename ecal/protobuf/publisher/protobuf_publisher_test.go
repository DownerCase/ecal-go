package publisher

import (
	"testing"

	"github.com/DownerCase/ecal-go/protos"
)

func TestProtobufPublisher(t *testing.T) {
	pub, err := New(&protos.Person{})

	if err != nil {
		t.Error(err)
	}
	defer pub.Delete()

	if err := pub.Create("testing"); err != nil {
		t.Error(err)
	}

	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

	person := &protos.Person{Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	// TODO: Check datatype information
	if err := pub.Send(person); err != nil {
		t.Error("Failed to send message", err)
	}
}
