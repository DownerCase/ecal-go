package protobuf

import (
	"testing"

	"github.com/DownerCase/ecal-go/protos"
)

func TestFullName(t *testing.T) {
	t.Parallel()

	expectedName := "pb.People.Person"
	if fn := GetFullName(&protos.Person{}); fn != expectedName {
		t.Error("Expected: ", expectedName, " Actual: ", fn)
	}
}
