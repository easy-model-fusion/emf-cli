package config

import (
	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/viper"
	"os"
	"testing"
)

type viperTestStructureOne struct {
	Name string
}

type viperTestStructureTwo struct {
	Name int
}

func TestGetViperConfig_Success(t *testing.T) {

	// Create file
	file, err := os.OpenFile("config.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	// Load the configuration file
	err = GetViperConfig()

	// Assert that the load method did not return an error
	test.AssertEqual(t, err, nil, "No error should have been raised")

	// Delete file
	file.Close()
	os.Remove(file.Name())
}

func TestGetViperConfig_Error(t *testing.T) {
	// Load the configuration file
	err := GetViperConfig()

	// Assert that the load method did return an error because no conf file in project
	test.AssertNotEqual(t, err, nil, "An error should have been raised")
}

func TestGetViperItem_Success(t *testing.T) {

	// Set up a test Viper configuration
	viper.Reset()
	testValue := []viperTestStructureOne{
		{Name: "name"},
	}
	viper.Set("test", testValue)

	// Call the GetViperItem function
	var result []viperTestStructureOne
	err := GetViperItem("test", &result)

	// Assert that the item was returned successfully
	test.AssertEqual(t, err, nil, "Error while retrieving the config item.")
	test.AssertEqual(t, len(result), len(testValue), "Expected the result item to be the same as the initial item.")
}

func TestGetViperItem_Error(t *testing.T) {

	// Set up a test Viper configuration
	viper.Reset()
	testValue := []viperTestStructureOne{
		{Name: "name"},
	}
	viper.Set("test", testValue)

	// Call the GetViperItem function with a non-existent key
	var result []viperTestStructureTwo
	err := GetViperItem("test", &result)

	// Assert
	test.AssertNotEqual(t, err, nil, "Error while retrieving the config item.")
}

func TestWriteViperConfig_Success(t *testing.T) {

	// Create file
	file, err := os.OpenFile("config.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	// Load the configuration file
	viper.Reset()
	viper.SetConfigFile(file.Name())
	err = WriteViperConfig()

	// Assert that the write method did not return an error
	test.AssertEqual(t, err, nil, "No error should have been raised")

	// Delete file
	file.Close()
	os.Remove(file.Name())
}

func TestWriteViperConfig_Error(t *testing.T) {
	// Load the configuration file
	viper.Reset()
	err := WriteViperConfig()

	// Assert that the load method did return an error because no conf file in project
	test.AssertNotEqual(t, err, nil, "An error should have been raised")
}
