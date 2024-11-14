package publisher

import (
	"testing"
)

func TestProtobufPublisher(t *testing.T) {
	pub, err := New()

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

	// TODO: Check datatype information
	if err := pub.Send("my message"); err != nil {
		t.Error("Failed to send message", err)
	}
}

