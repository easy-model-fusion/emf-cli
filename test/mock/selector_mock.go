package mock

import (
	"github.com/easy-model-fusion/emf-cli/internal/model"
)

type MockModelSelector struct {
	SelectorModel   model.Model
	SelectorError   error
	SelectorWarning string
}

func (d *MockModelSelector) SelectTransformerModel(models model.Models, configModelsMap map[string]model.Model) (model.Model, string, error) {

	return d.SelectorModel, d.SelectorWarning, d.SelectorError
}
