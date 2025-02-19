package testutil

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func NewBinarySubscriber(t *testing.T, topic string) *ecal.BinarySubscriber {
	t.Helper()

	sub, err := ecal.NewBinarySubscriber(topic)
	if err != nil {
		t.Error(err)
	}

	return sub
}

func NewStringSubscriber(t *testing.T, topic string) *ecal.StringSubscriber {
	t.Helper()

	sub, err := ecal.NewStringSubscriber(topic)
	if err != nil {
		t.Error(err)
	}

	return sub
}

func NewProtobufSubscriber[U any, T ecal.ProtoMessage[U]](t *testing.T, topic string) *ecal.ProtobufSubscriber[T] {
	t.Helper()

	sub, err := ecal.NewProtobufSubscriber[U, T](topic)
	if err != nil {
		t.Error(err)
	}

	return sub
}
