package app

import (
	"github.com/easy-model-fusion/client/test"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	initLogger()
}

func TestLoadNotExistentConfFile(t *testing.T) {
	// Load the configuration file
	err := LoadConfFile(".")
	// Assert that the load method did return an error because no conf file in project
	test.AssertNotEqual(t, err, nil)
}

func TestLoadConfFile(t *testing.T) {
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(dname)
	// Create a temporary config file for the test
	file, err := os.Create(filepath.Join(dname, "config.yaml"))
	defer file.Close()
	test.AssertEqual(t, err, nil, "Error while creating conf file.")

	// Write some content to the config file
	_, err = file.WriteString("key: value")
	test.AssertEqual(t, err, nil, "Error while writing into conf file.")

	// Load the configuration file
	err = LoadConfFile(dname)

	// Assert that the load method did not return any error
	test.AssertEqual(t, err, nil, "Error while loading conf file.")

	// Assert that the loaded configuration has the expected values
	test.AssertEqual(t, "value", viper.GetString("key"))
}
