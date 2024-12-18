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

	initResult := ecal.Initialize(ecal.NewConfig(), "go_test", ecal.CDefault)
	if initResult == 1 {
		t.Fatal("eCAL already initialized")
	} else if initResult != 0 {
		t.Fatalf("eCAL failed to initialize with error %v", initResult)
	}

	// Test double initialization
	secondInit := ecal.Initialize(ecal.NewConfig(), "go_test2", ecal.CPublisher)
	if secondInit != 1 {
		t.Errorf("Second initialize returned %v", secondInit)
	}

	if !ecal.IsInitialized() {
		t.Error("IsInitialized return false, expected true")
	}

	if !ecal.IsComponentInitialized(ecal.CPublisher) {
		t.Error("Expected publisheCPublisher to be initialised")
	}

	if !ecal.SetUnitName("go_test_set_name") {
		t.Error("Failed to set unit name")
	}

	if !ecal.Ok() {
		t.Error("eCAL not Ok")
	}

	finalizeResult := ecal.Finalize()
	if finalizeResult != 0 {
		t.Errorf("Failed to finalize with error %v", finalizeResult)
	}

	secondFinalize := ecal.Finalize()
	// We've called Initialize twice so 2 calls to Finalize are needed
	if secondFinalize != 0 {
		t.Errorf("Second finalize returned %v", secondFinalize)
	}

	thirdFinalize := ecal.Finalize()
	if thirdFinalize != 1 {
		t.Errorf("Expected Finalize to already be done, recevied %v", thirdFinalize)
	}

	if ecal.Ok() {
		t.Error("eCAL Ok after being finalized")
	}
}
