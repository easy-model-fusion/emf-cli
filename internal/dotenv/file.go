package dotenv

import (
	"github.com/joho/godotenv"
)

// GetEnvValue returns the value of a given environment variable
func GetEnvValue(key string) (value string, err error) {
	env, err := godotenv.Read()
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
	env, err := godotenv.Read()
	if err != nil {
		return err
	}
	env[key] = value
	return godotenv.Write(env, "./.env")
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
