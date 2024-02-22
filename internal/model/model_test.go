package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/codegen"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestModel_GenClass(t *testing.T) {
	model := Model{
		Name: "stabilityai/sdxl-turbo",
		Path: "build/stabilityai/sdxl-turbo",
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

func TestModel_GenFile(t *testing.T) {
	model := Model{
		Name:        "stabilityai/sdxl-turbo",
		Path:        "build/stabilityai/sdxl-turbo",
		PipelineTag: TextToImage,
	}

	file := model.GenFile()
	gen := codegen.NewPythonCodeGenerator(true)
	result, err := gen.Generate(file)

	t.Logf("\n%s", result)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, file.Name, "StabilityaiSdxlTurbo.py", "The file name should be formatted correctly.")
	test.AssertEqual(t, len(file.Classes), 1, "The file should contain one class.")
	test.AssertEqual(t, len(file.HeaderComments), 2, "The file should contain two header comments.")
}

func TestModel_GetFormattedModelName(t *testing.T) {
	model := Model{
		Name: "stabilityai/sdxl-turbo",
	}

	test.AssertEqual(t, model.GetFormattedModelName(), "StabilityaiSdxlTurbo", "The model name should be formatted correctly.")
}

func TestModel_GetPipelineTagAbstractClassName(t *testing.T) {
	model := Model{
		PipelineTag: TextToImage,
	}

	test.AssertEqual(t, model.GetPipelineTagAbstractClassName(), "ModelTextToImage", "The model name should be formatted correctly.")

	model.PipelineTag = TextGeneration

	test.AssertEqual(t, model.GetPipelineTagAbstractClassName(), "ModelTextToText", "The model name should be formatted correctly.")

	model.PipelineTag = "unknown"
	test.AssertEqual(t, model.GetPipelineTagAbstractClassName(), "", "The model name should be formatted correctly.")
}
