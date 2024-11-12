package ecal

import (
	"testing"
)

func TestNewPublishers(t *testing.T) {
	for range 100 {
		ptr, err := NewPublisher()
		if err != nil {
			t.Error(err)
		}
		defer func() {
			if !DestroyPublisher(&ptr) {
				t.Error("Failed to cleanup!")
			}
		}()
		t.Log(ptr)
	}
}

func TestPublisher(t *testing.T) {
	pub, err := NewPublisher()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if !DestroyPublisher(&pub) {
			t.Error("Failed to cleanup!")
		}
		if DestroyPublisher(&pub) {
			t.Error("Destroyed publisher twice!")
		}
	}()

	if err := pub.Create("testing"); err != nil {
		t.Error(err)
	}
	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

}
