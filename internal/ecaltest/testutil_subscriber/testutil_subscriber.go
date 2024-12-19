package testutilsubscriber

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/subscriber"
)

func NewGenericSubscriber(t *testing.T, topic string) *subscriber.Subscriber {
	t.Helper()

	sub, err := subscriber.New(topic, subscriber.DataType{})
	if err != nil {
		t.Error(err)
	}

	return sub
}
