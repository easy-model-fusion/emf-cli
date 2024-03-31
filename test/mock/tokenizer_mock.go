package mock

import (
	"github.com/easy-model-fusion/emf-cli/internal/model"
)

type MockModels struct {
	GetModelList model.Models
	GetError     error
}

func (m *MockModels) GetModels() (model.Models, error) {
	return m.GetModelList, m.GetError
}
