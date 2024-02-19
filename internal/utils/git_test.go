package utils

import (
	"github.com/easy-model-fusion/client/test"
	"os"
	"path/filepath"
	"testing"
)

func TestCloneSDK(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the SDK
	err = CloneSDK("v1.0.0", tmpDir)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Check if the SDK was cloned (.git folder exists)
	_, err = os.Stat(filepath.Join(tmpDir, ".git"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestGetLatestTag(t *testing.T) {
	// Test with a valid tag
	tag, err := GetLatestTag("sdk")
	test.AssertEqual(t, err, nil, "Expected no error")
	t.Log(tag)

	// Test with an invalid tag
	_, err = GetLatestTag("invalid")
	test.AssertNotEqual(t, err, nil, "Expected error")
}
