package dotenv

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// Test all dotenv files methods
func TestFullyDotEnv(t *testing.T) {
	// Create full test suite with a .env file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Add new variable
	err := AddNewEnvVariable("key", "value")
	test.AssertEqual(t, err, nil, "No error expected while adding environment variable")

	// Verify if the variable exist
	exist, err := EnvVariableExists("key")
	test.AssertEqual(t, err, nil, "No error expected while searching for environment variable")
	test.AssertEqual(t, exist, true, "key variable should have been created")

	// Get the variable value
	value, err := GetEnvValue("key")
	test.AssertEqual(t, err, nil, "No error expected while fetching environment variable")
	test.AssertEqual(t, value, "value", "Value of variable should be equal to the one entered earlier")

	// Set new environment variable with same key
	newKey, err := SetNewEnvKey("key")
	test.AssertEqual(t, err, nil, "No error expected while creating unique key for environment variable")
	test.AssertEqual(t, newKey, "key_2", "Value of key should be unique")

	// Remove variable
	err = RemoveEnvVariable("key")
	test.AssertEqual(t, err, nil, "No error expected while removing environment variable")

	// Verify if the variable removed
	exist, err = EnvVariableExists("key")
	test.AssertEqual(t, err, nil, "No error expected while searching for invalid environment variable")
	test.AssertEqual(t, exist, false, "key variable should have been removed")
}

// Tests GetEnvValue with no .env file
func TestGetEnvValue_WithNoEnvFile(t *testing.T) {
	// Get the variable value
	value, err := GetEnvValue("key")
	test.AssertEqual(t, err, nil, "No error expected while fetching environment variable")
	test.AssertEqual(t, value, "", "No variable with given key should be found")
}

// Tests GetEnvValue with invalid variable key
func TestRemoveEnvVariable_WithInvalidKey(t *testing.T) {
	// Create full test suite with a .env file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Remove variable
	err := RemoveEnvVariable("invalid key")
	test.AssertEqual(t, err.Error(), "environment variable does not exist", "Error expected while removing not existing environment variable")
}
