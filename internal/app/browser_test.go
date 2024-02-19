package app

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"testing"
)

// TestGetDownloadedModels tests GetDownloadedModels method
func TestGetDownloadedModels(t *testing.T) {
	test.CreateModelsFolderFullTestSuite(t)
	defer os.RemoveAll(ModelsDownloadPath)

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
