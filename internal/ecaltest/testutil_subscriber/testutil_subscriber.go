package testutilsubscriber

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/subscriber"
)

func NewGenericSubscriber(t *testing.T, topic string) *subscriber.Subscriber {
	t.Helper()

	sub, err := subscriber.New()
	if err != nil {
		t.Error(err)
	}

	if err := sub.Create(topic, subscriber.DataType{}); err != nil {
		t.Error(err)
	}

	return sub
}
