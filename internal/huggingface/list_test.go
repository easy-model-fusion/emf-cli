package huggingface

import (
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/easy-model-fusion/emf-cli/test"
	"reflect"
	"testing"
)

func TestGetModels(t *testing.T) {
	h := NewHuggingFace(BaseUrl, "")
	models, err := h.GetModels(string(model.TextToImage), 10)
	test.AssertEqual(t, err, nil, "The api call should've passed.")

	for _, apiModel := range models {
		valid := utils.SliceContainsItem(model.AllModules, apiModel.Config.Module)
		test.AssertEqual(t, valid, true, "Model's module should be valid.")
	}
}

func Test_getModelsByModules(t *testing.T) {
	validModels := []model.Model{
		{
			Name: "model1",
			Config: model.Config{
				Module: model.DIFFUSERS,
			},
		},
		{
			Name: "model2",
			Config: model.Config{
				Module: model.TRANSFORMERS,
			},
		},
		{
			Name: "model3",
			Config: model.Config{
				Module: model.DIFFUSERS,
			},
		},
	}
	invalidModels := []model.Model{
		{
			Name: "model1",
			Config: model.Config{
				Module: model.DIFFUSERS,
			},
		},
		{
			Name: "model2",
			Config: model.Config{
				Module: model.TRANSFORMERS,
			},
		},
		{
			Name: "model3",
			Config: model.Config{
				Module: model.DIFFUSERS,
			},
		},
		{
			Name: "invalid_model",
			Config: model.Config{
				Module: "INVALID",
			},
		},
	}
	type args struct {
		models []model.Model
	}
	tests := []struct {
		name               string
		args               args
		wantReturnedModels []model.Model
	}{
		{name: "All valid models", args: args{validModels}, wantReturnedModels: validModels},
		{name: "Not all valid models", args: args{invalidModels}, wantReturnedModels: validModels},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReturnedModels := getModelsByModules(tt.args.models); !reflect.DeepEqual(gotReturnedModels, tt.wantReturnedModels) {
				t.Errorf("getModelsByModules() = %v, want %v", gotReturnedModels, tt.wantReturnedModels)
			}
		})
	}
}
