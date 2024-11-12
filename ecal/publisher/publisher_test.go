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
		defer func() {
			if !Destroy(&ptr) {
				t.Error("Failed to cleanup!")
			}
		}()
		t.Log(ptr)
	}
}

func TestPublisher(t *testing.T) {
	pub, err := New()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if !Destroy(&pub) {
			t.Error("Failed to cleanup!")
		}
		if Destroy(&pub) {
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
