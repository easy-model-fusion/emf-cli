package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/codegen"
	"testing"
)

func TestModel_GenClass(t *testing.T) {
	model := Model{
		Name: "stabilityai/sdxl-turbo",
		Config: Config{
			Path: "build/stabilityai/sdxl-turbo",
		},
	}

	class := model.GenClass()
	gen := codegen.NewPythonCodeGenerator(true)
	result, err := gen.Generate(&codegen.File{
		Name:    "test",
		Classes: []*codegen.Class{class},
	})

	t.Logf("\n%s", result)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}
