package ecal

import (
	"testing"
)

func TestVersionString(t *testing.T) {
	t.Parallel()
	version := GetVersionString()
	if version == "" {
		t.Error("GetVersionString returned empty string")
	}
	t.Log(version)
}

func TestVersionDateString(t *testing.T) {
	t.Parallel()
	buildDate := GetVersionDateString()
	if buildDate == "" {
		t.Error("GetVersionDateString returned empty string")
	}
	t.Log(buildDate)
}

func TestGetVersion(t *testing.T) {
	t.Parallel()
	version := GetVersion()
	t.Log(version)
}

func TestInitializeFinalize(t *testing.T) {
	if IsInitialized() {
		t.Error("eCAL pre-initialized...")
	}

	initResult := Initialize(NewConfig(), "go_test", CDefault)
	if initResult == 1 {
		t.Fatal("eCAL already initialized")
	} else if initResult != 0 {
		t.Fatalf("eCAL failed to initialize with error %v", initResult)
	}

	// Test double initialization
	secondInit := Initialize(NewConfig(), "go_test2", CPublisher)
	if secondInit != 1 {
		t.Errorf("Second initialize returned %v", secondInit)
	}

	if !IsInitialized() {
		t.Error("IsInitialized return false, expected true")
	}
	if !IsComponentInitialized(CPublisher) {
		t.Error("Expected publisheCPublisher to be initialised")
	}

	if !SetUnitName("go_test_set_name") {
		t.Error("Failed to set unit name")
	}

	if !Ok() {
		t.Error("eCAL not Ok")
	}

	finalizeResult := Finalize()
	if finalizeResult != 0 {
		t.Errorf("Failed to finalize with error %v", finalizeResult)
	}

	secondFinalize := Finalize()
	// We've called Initialize twice so 2 calls to Finalize are needed
	if secondFinalize != 0 {
		t.Errorf("Second finalize returned %v", secondFinalize)
	}

	thirdFinalize := Finalize()
	if thirdFinalize != 1 {
		t.Errorf("Expected Finalize to already be done, recevied %v", thirdFinalize)
	}

	if Ok() {
		t.Error("eCAL Ok after being finalized")
	}
}
