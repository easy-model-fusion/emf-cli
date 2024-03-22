package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/codegen"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestModel_GenClass(t *testing.T) {
	model := Model{
		Name:        "stabilityai/sdxl-turbo",
		Path:        "build/stabilityai/sdxl-turbo",
		PipelineTag: huggingface.TextToImage,
		Module:      huggingface.DIFFUSERS,
		Class:       "DiffusionPipeline",
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

func TestModel_DiffuserGenFile(t *testing.T) {
	model := Model{
		Name:        "stabilityai/sdxl-turbo",
		Path:        "build/stabilityai/sdxl-turbo",
		PipelineTag: huggingface.TextToImage,
		Module:      huggingface.DIFFUSERS,
		Class:       "DiffusionPipeline",
	}

	file := model.GenFile()
	gen := codegen.NewPythonCodeGenerator(true)
	result, err := gen.Generate(file)

	t.Logf("\n%s", result)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, file.Name, "ModelStabilityaiSdxlTurbo.py", "The file name should be formatted correctly.")
	test.AssertEqual(t, len(file.Classes), 1, "The file should contain one class.")
	test.AssertEqual(t, len(file.HeaderComments), 2, "The file should contain two header comments.")
}

func TestModel_TransformersGenFile(t *testing.T) {
	model := Model{
		Name:        "microsoft/phi-2",
		Path:        "build/microsoft/phi-2/model",
		PipelineTag: huggingface.TextGeneration,
		Module:      huggingface.TRANSFORMERS,
		Class:       "AutoModelForCausalLM",
		Tokenizers: Tokenizers{
			{
				Class: "AutoTokenizer",
				Path:  "build/microsoft/phi-2/AutoTokenizer",
			},
		},
	}

	file := model.GenFile()
	gen := codegen.NewPythonCodeGenerator(true)
	result, err := gen.Generate(file)

	t.Logf("\n%s", result)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, file.Name, "ModelMicrosoftPhi2.py", "The file name should be formatted correctly.")
	test.AssertEqual(t, len(file.Classes), 1, "The file should contain one class.")
	test.AssertEqual(t, len(file.HeaderComments), 2, "The file should contain two header comments.")
}

func TestModel_GetFormattedModelName(t *testing.T) {
	model := Model{
		Name:        "stabilityai/sdxl-turbo",
		Path:        "build/stabilityai/sdxl-turbo",
		PipelineTag: huggingface.TextToImage,
		Module:      huggingface.DIFFUSERS,
	}

	test.AssertEqual(t, model.GetFormattedModelName(), "ModelStabilityaiSdxlTurbo", "The model name should be formatted correctly.")
}

func TestModel_GetPipelineTagAbstractClassName(t *testing.T) {
	model := Model{
		Module: huggingface.DIFFUSERS,
	}

	test.AssertEqual(t, model.GetSDKClassNameWithModule(), "ModelDiffusers", "The model name should be formatted correctly.")

	model.Module = huggingface.TRANSFORMERS

	test.AssertEqual(t, model.GetSDKClassNameWithModule(), "ModelTransformers", "The model name should be formatted correctly.")

	model.Module = "unknown"
	test.AssertEqual(t, model.GetSDKClassNameWithModule(), "", "The model name should be formatted correctly.")
}
