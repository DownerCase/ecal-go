package logging

import (
	"os"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	_ "github.com/DownerCase/ecal-go/ecal/types"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
)

func expectMessageIsFromHost(t *testing.T, msg LogMessage) {
	host, err := os.Hostname()
	if err != nil && msg.Host != host {
		t.Error("Host mismatch", host, "!=", msg.Host)
	}
	if int(msg.Pid) != os.Getpid() {
		t.Error("Mismatch pid", os.Getpid(), "!=", msg.Pid)
	}
}

func receiveMessage(t *testing.T, msg string, level Level) bool {
	logs := GetLogging()

	for _, rmsg := range logs.Messages {
		if rmsg.Content == msg {
			expectMessageIsFromHost(t, rmsg)
			if rmsg.Level != level {
				t.Error("Mismatch log level", rmsg.Level, "!=", level)
			}
			return true
		}
	}
	return false
}

func TestGetLogging(t *testing.T) {
	// Given: eCAL Initialized
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	// When: Logging a message
	testMessage := "This is a test log message"
	level := LevelError
	Log(level, testMessage)

	// Expect: To receieve that message
	time.Sleep(5 * time.Millisecond)
	if !receiveMessage(t, testMessage, level) {
		t.Error("Could not find expected message:", testMessage)
	}
	// Expect: To not be able to receive it again
	if receiveMessage(t, testMessage, level) {
		t.Error("Recevied duplicate message:", testMessage)
	}
}
