package publisher_test

import (
	"testing"

	testutilpublisher "github.com/DownerCase/ecal-go/internal/ecaltest/string/testutil_publisher"
)

func TestStringPublisher(t *testing.T) {
	pub := testutilpublisher.NewStringPublisher(t, "test_string_publisher")
	defer pub.Delete()

	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

	// TODO: Check datatype information
	if err := pub.Send("my message"); err != nil {
		t.Error("Failed to send message", err)
	}
}
