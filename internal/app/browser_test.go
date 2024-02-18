package app

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetDownloadedModels tests GetDownloadedModels method
func TestGetDownloadedModels(t *testing.T) {
	err := os.Mkdir(ModelsDownloadPath, 0755)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(ModelsDownloadPath)

	setupTestDir(t, ModelsDownloadPath)
	models, err := GetDownloadedModelNames()
	if err != nil {
		t.Errorf("Error getting downloaded models: %v", err)
	}

	expected := []string{"model1/weights", "model2/weights", "model3"}
	// Tests if the method returned the expected number of models
	if len(models) != len(expected) {
		t.Errorf("Expected %d models, but got %d", len(expected), len(models))
	}
	// Tests if the method returned the expected model names
	for i, model := range models {
		if model != expected[i] {
			t.Errorf("Expected model %s, but got %s", expected[i], model)
		}
	}
}

// setupTestDir setups a test dir with different mock models
func setupTestDir(t *testing.T, dir string) string {
	err := os.MkdirAll(filepath.Join(dir, "model1", "weights"), 0755)
	if err != nil {
		t.Fatalf("Error creating test directory: %v", err)
	}
	err = os.MkdirAll(filepath.Join(dir, "model2", "weights"), 0755)
	if err != nil {
		t.Fatalf("Error creating test directory: %v", err)
	}
	err = os.MkdirAll(filepath.Join(dir, "model3"), 0755)
	if err != nil {
		t.Fatalf("Error creating test directory: %v", err)
	}

	return dir
}
