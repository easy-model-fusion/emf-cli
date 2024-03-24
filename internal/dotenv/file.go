package dotenv

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/joho/godotenv"
)

// GetEnvValue returns the value of a given environment variable
func GetEnvValue(key string) (value string, err error) {
	// Check if the .env file exists
	exist, err := fileutil.IsExistingPath(".env")
	if err != nil || !exist {
		return "", err
	}

	// Find environment variable
	env, err := godotenv.Read(".env")
	if err == nil {
		value = env[key]
	}
	return value, err
}

// EnvVariableExists returns true if an environment variable with the given key exists
func EnvVariableExists(key string) (bool, error) {
	value, err := GetEnvValue(key)
	return value != "", err
}

// AddNewEnvVariable adds a new environment variable
func AddNewEnvVariable(key string, value string) error {
	env, err := godotenv.Read(".env")
	if err != nil {
		return err
	}
	env[key] = value
	return godotenv.Write(env, ".env")
}

// RemoveEnvVariable removes an environment variable from a .env file
func RemoveEnvVariable(key string) error {
	// Read the current environment variables from .env file
	env, err := godotenv.Read(".env")
	if err != nil {
		return err
	}

	// Check if the variable exists in the environment
	if _, exists := env[key]; !exists {
		return errors.New("environment variable does not exist")
	}

	// Remove the variable from the environment
	delete(env, key)

	// Write the updated environment back to .env file
	return godotenv.Write(env, ".env")
}

// SetNewEnvKey sets new unique env key
func SetNewEnvKey(key string) (string, error) {
	exist, err := EnvVariableExists(key)
	if err != nil {
		return "", err
	}
	if exist {
		key += "_2"
		return SetNewEnvKey(key)
	}

	return key, nil
}
