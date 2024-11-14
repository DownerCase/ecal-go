package publisher

import (
	"testing"
)

func TestNewPublishers(t *testing.T) {
	for range 100 {
		ptr, err := New()
		if err != nil {
			t.Error(err)
		}
		defer ptr.Delete()
		t.Log(ptr)
	}
}

func TestPublisher(t *testing.T) {
	pub, err := New()
	if err != nil {
		t.Error(err)
	}
	defer pub.Delete()

	if err := pub.Create("testing", DataType{}); err != nil {
		t.Error(err)
	}
	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

}
