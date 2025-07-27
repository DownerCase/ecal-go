package ecal_test

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func TestVersionString(t *testing.T) {
	t.Parallel()

	version := ecal.GetVersionString()
	if version == "" {
		t.Error("GetVersionString returned empty string")
	}

	t.Log(version)
}

func TestVersionDateString(t *testing.T) {
	t.Parallel()

	buildDate := ecal.GetVersionDateString()
	if buildDate == "" {
		t.Error("GetVersionDateString returned empty string")
	}

	t.Log(buildDate)
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	version := ecal.GetVersion()
	t.Log(version)
}

func TestInitializeFinalize(t *testing.T) {
	if ecal.IsInitialized() {
		t.Error("eCAL pre-initialized...")
	}

	if !ecal.Initialize(ecal.NewConfig(), "go_test", ecal.CDefault) {
		t.Fatalf("eCAL failed to initialize")
	}

	// Test double initialization
	if ecal.Initialize(ecal.NewConfig(), "go_test2", ecal.CPublisher) {
		t.Errorf("Second initialize returned true")
	}

	if !ecal.IsInitialized() {
		t.Error("IsInitialized return false, expected true")
	}

	if !ecal.IsComponentInitialized(ecal.CPublisher) {
		t.Error("Expected publisheCPublisher to be initialised")
	}

	if !ecal.Ok() {
		t.Error("eCAL not Ok")
	}

	if !ecal.Finalize() {
		t.Errorf("Failed to finalize")
	}

	// eCAL does not reference count (as of v6)
	if ecal.Finalize() {
		t.Errorf("Expected second finalize to fail")
	}

	if ecal.Ok() {
		t.Error("eCAL Ok after being finalized")
	}
}
