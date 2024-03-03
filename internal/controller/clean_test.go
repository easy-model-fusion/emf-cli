package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"testing"
)

func TestRunClean(t *testing.T) {
	app.SetUI(test.NewMockUI())
	app.SetGit(test.NewMockGit())

	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// dry run
	RunClean(false, false)

	// now create a file in the build directory
	err := os.Mkdir(cleanDirName, os.ModePerm)
	test.AssertEqual(t, err, nil, "Error creating build directory")

	// run clean
	RunClean(false, false)

	// check if the build directory is deleted
	_, err = os.Stat(cleanDirName)
	test.AssertEqual(t, os.IsNotExist(err), true, "Build directory not deleted")

	// add 3 models to config
	err = config.AddModels([]model.Model{
		{
			Name: "model1",
			Path: "models/model1",
		},
		{
			Name: "model2",
			Path: "models/model2",
		},
		{
			Name: "model3",
			Path: "models/model3",
		},
	})
	test.AssertEqual(t, err, nil, "Error adding models to config")

	// test allFlagDelete
	RunClean(true, false)

	// count files in models directory, should be 3 (no confirmation)
	files, err := os.ReadDir("models")
	test.AssertEqual(t, err, nil, "Error reading models directory")
	test.AssertEqual(t, len(files), test.FullTestSuiteModelsCount, "Models directory should not be empty")

	// test allFlagDelete with confirmation
	RunClean(true, true)

	// check if the models directory is deleted (with confirmation)
	_, err = os.Stat("models")
	test.AssertEqual(t, os.IsNotExist(err), true, "Models directory not deleted")
}
