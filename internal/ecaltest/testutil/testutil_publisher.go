package testutil

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func NewBinaryPublisher(t *testing.T, topic string) *ecal.BinaryPublisher {
	t.Helper()

	pub, err := ecal.NewBinaryPublisher(topic)
	if err != nil {
		t.Error(err)
	}

	return pub
}

func NewStringPublisher(t *testing.T, topic string) *ecal.StringPublisher {
	t.Helper()

	pub, err := ecal.NewStringPublisher(topic)
	if err != nil {
		t.Error(err)
	}

	return pub
}

func NewProtobufPublisher[U any, T ecal.ProtoMessage[U]](t *testing.T, topic string) *ecal.ProtobufPublisher[T] {
	t.Helper()

	pub, err := ecal.NewProtobufPublisher[U, T](topic)
	if err != nil {
		t.Error(err)
	}

	return pub
}
