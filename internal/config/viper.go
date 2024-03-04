package config

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// GetViperConfig Config loaded and return an error upon failure
func GetViperConfig(confDirPath string) (err error) {
	count := 0
	for count < 3 {
		err = Load(confDirPath)
		if err != nil {
			count++
			confDirPath = UpdateConfigFilePath()
		} else {
			return err
		}
	}

	pterm.Error.Println(fmt.Sprintf("Error loading config file after %d attempts: %s", count, err))
	return err
}

// GetViperItem Store the key data into the target
func GetViperItem(key string, target interface{}) (err error) {
	if err = viper.UnmarshalKey(key, target); err != nil {
		pterm.Error.Println(fmt.Sprintf("Error reading config file : %s", err))
		return err
	}
	return nil
}

// WriteViperConfig Attempt to write the configuration file
func WriteViperConfig() (err error) {
	if err = viper.WriteConfig(); err != nil {
		pterm.Error.Println(fmt.Sprintf("Error writing to config file : %s", err))
		return err
	}
	return nil
}
