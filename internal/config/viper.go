package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// GetViperConfig Config loaded and return an error upon failure
func GetViperConfig(confDirPath string) (err error) {
	count := 0
	// Store the current config file path to restore it in case of failure
	tempConfPath := FilePath

	// Try to load the config file 3 times max
	for count < 3 {
		err = Load(confDirPath)
		if err != nil {
			count++
			confDirPath = UpdateConfigFilePath()
		} else {
			return nil // Success
		}
	}

	// restore the original config file path
	FilePath = tempConfPath
	return fmt.Errorf("error loading config file after %d attempts: %s", count, err)
}

// GetViperItem Store the key data into the target
func GetViperItem(key string, target interface{}) (err error) {
	if err = viper.UnmarshalKey(key, target); err != nil {
		return fmt.Errorf("error reading config file : %s", err)
	}
	return nil
}

// WriteViperConfig Attempt to write the configuration file
func WriteViperConfig() (err error) {
	if err = viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config file : %s", err)
	}
	return nil
}
