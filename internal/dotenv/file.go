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

// AddNewEnvVariable adds a new environment variable
func AddNewEnvVariable(key string, value string) error {
	env, err := godotenv.Unmarshal(key + "=" + value)
	if err != nil {
		return err
	}
	return godotenv.Write(env, "./.env")
}
