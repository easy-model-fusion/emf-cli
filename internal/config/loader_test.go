package config

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	app.Init("", "")
}

func TestLoadNotExistentConfFile(t *testing.T) {
	// Load the configuration file
	err := Load(".")
	// Assert that the load method did return an error because no conf file in project
	test.AssertNotEqual(t, err, nil)
}

func TestLoad(t *testing.T) {
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(dname)
	// Create a temporary config file for the test
	file, err := os.Create(fileutil.PathJoin(dname, "config.yaml"))
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			t.Error(err)
		}
	}(file)
	test.AssertEqual(t, err, nil, "Error while creating conf file.")

	// Write some content to the config file
	_, err = file.WriteString("key: value")
	test.AssertEqual(t, err, nil, "Error while writing into conf file.")

	// Load the configuration file
	err = Load(dname)

	// Assert that the load method did not return any error
	test.AssertEqual(t, err, nil, "Error while loading conf file.")

	// Assert that the loaded configuration has the expected values
	test.AssertEqual(t, "value", viper.GetString("key"))
}

// Tests UpdateConfigFilePath
func TestUpdateConfigFilePath(t *testing.T) {
	//create mock UI
	ui := mock.MockUI{UserInputResult: "path/test"}
	app.SetUI(ui)

	// Update the configuration file path
	path := UpdateConfigFilePath()

	//Assertions
	test.AssertEqual(t, "path/test", path)
}
