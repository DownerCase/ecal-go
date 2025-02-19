package ecal_test

import (
	"testing"

	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil"
	"github.com/DownerCase/ecal-go/protos"
)

func TestProtobufPublisher(t *testing.T) {
	pub := testutil.NewProtobufPublisher[protos.Person](t, "testing_protobuf_publisher")
	defer pub.Delete()

	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

	person := &protos.Person{
		Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	// TODO: Check datatype information
	pub.Send(person)
}
