package config

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"testing"
)

// Define test structures for use in the tests
type viperTestStructureOne struct {
	Name string
}
type viperTestStructureTwo struct {
	Name int
}

// TestGetViperConfig_Success tests the successful loading of the Viper configuration.
func TestGetViperConfig_Success(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	FilePath = "."
	// Load the configuration file
	err := GetViperConfig(FilePath)

	// Assert that the load method did not return an error
	test.AssertEqual(t, err, nil, "No error should have been raised")
}

// TestGetViperItem_Success tests the successful retrieval of an item from the Viper configuration.
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

// TestGetViperItem_Error tests the case where there is an error retrieving an item from the Viper configuration.
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

	// Assert that an error was raised
	test.AssertNotEqual(t, err, nil, "Error while retrieving the config item.")
}

// TestWriteViperConfig_Success tests the successful writing of the Viper configuration.
func TestWriteViperConfig_Success(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Load the configuration file
	viper.Reset()
	viper.SetConfigFile("config.yaml")
	err := WriteViperConfig()

	// Assert that the write method did not return an error
	test.AssertEqual(t, err, nil, "No error should have been raised")
}

// TestWriteViperConfig_Error tests the case where there is an error writing the Viper configuration.
func TestWriteViperConfig_Error(t *testing.T) {
	// Load the configuration file
	viper.Reset()
	err := WriteViperConfig()

	// Assert that the load method did return an error because no conf file in project
	test.AssertNotEqual(t, err, nil, "An error should have been raised")
}
