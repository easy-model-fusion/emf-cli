package tokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

// TestSelectTransformerModel_Success tests successfull selection
func TestSelectTransformerModel_Success(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.DIFFUSERS,
	})
	models = append(models, model.Model{
		Name:   "model2",
		Module: huggingface.TRANSFORMERS,
	})
	models = append(models, model.Model{
		Name:   "model3",
		Module: huggingface.TRANSFORMERS,
	})
	// Create ui mock
	ui := mock.MockUI{SelectResult: "model2"}
	app.SetUI(ui)

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	configModelsMap := models.Map()

	ic := SelectModelController{}

	selectedModel := ic.SelectTransformerModel(models, configModelsMap)

	// Assert that the selected model is as expected
	if selectedModel.Name != "model2" {
		t.Errorf("Expected selected model name to be 'model2', got '%s'",
			selectedModel.Name)
	}
}
