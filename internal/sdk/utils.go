package sdk

import (
	"errors"
	"fmt"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// checkForUpdates Check for updates and return a whether there is an update or not
// Compares the current version of the SDK with the latest tag on GitHub
// Returns the latest tag and a boolean indicating if there is an update
func checkForUpdates() (string, bool) {
	version := viper.GetString("sdk-tag")
	if version == "" {
		return "", false
	}

	tag, err := utils.GetLatestTag("sdk")
	if err != nil {
		return "", false
	}

	return tag, tag != version
}

// canSendUpdateSuggestion Check if the user has already been suggested to update
func canSendUpdateSuggestion() bool {
	if !viper.IsSet("update-suggested") {
		return true
	}
	return !viper.GetBool("update-suggested")
}

// ResetUpdateSuggestion Reset the update suggestion
func ResetUpdateSuggestion() {
	setUpdateSuggestion(false)
}

// setUpdateSuggestion Set the update suggestion
func setUpdateSuggestion(value bool) {
	viper.Set("update-suggested", value)
	_ = viper.WriteConfig() // ignore error
}

// SendUpdateSuggestion Send an update suggestion to the user, if there is an update available and if they haven't been suggested before
// This is used to avoid spamming the user with update suggestions
// The update suggestion is reset when the user updates the SDK
func SendUpdateSuggestion() {
	if !canSendUpdateSuggestion() {
		return
	}
	tag, ok := checkForUpdates()
	if !ok {
		return
	}

	pterm.DefaultBox.Println(fmt.Sprintf("A new version of the SDK (%s) is available!\n"+
		"To update, run 'emf-cli update'", tag))
	pterm.Println()

	setUpdateSuggestion(true)
}

// Upgrade the SDK to the latest version
// Make sure the config is loaded before calling this function
func Upgrade() error {
	tag, ok := checkForUpdates()
	if !ok {
		pterm.Info.Println("SDK is already up to date")
		return errors.New("sdk is already up to date")
	}

	// remove sdk folder
	spinner, _ := pterm.DefaultSpinner.Start("Cleaning up sdk folder...")
	err := os.RemoveAll("sdk")
	if err != nil {
		spinner.Fail(err)
		return err
	}
	// create sdk folder
	err = os.Mkdir("sdk", os.ModePerm)
	if err != nil {
		spinner.Fail(err)
		return err
	}
	spinner.Success()

	// clone sdk
	spinner, _ = pterm.DefaultSpinner.Start("Cloning latest sdk...")
	err = utils.CloneSDK(tag, filepath.Join("sdk"))
	if err != nil {
		spinner.Fail(err)
		return err
	}

	// update sdk tag
	viper.Set("sdk-tag", tag)

	err = viper.WriteConfig()
	if err != nil {
		spinner.Fail(err)
		return err
	}
	spinner.Success()

	pterm.Success.Println(fmt.Sprintf("SDK successfully updated to version %s", tag))
	return nil
}
