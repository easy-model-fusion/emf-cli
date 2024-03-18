package sdk

import (
	"errors"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
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

	tag, err := app.G().GetLatestTag("sdk")
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
		"To update, run 'emf-cli upgrade'", tag))
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
	spinner := app.UI().StartSpinner("Cleaning up sdk folder...")
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
	spinner = app.UI().StartSpinner("Cloning latest sdk...")
	err = app.G().CloneSDK(tag, "sdk")
	if err != nil {
		spinner.Fail(err)
		return err
	}
	spinner.Success()

	// Move files from sdk/sdk to sdk/
	spinner = app.UI().StartSpinner("Reorganizing SDK files")
	err = fileutil.MoveFiles(filepath.Join("sdk", "sdk"), "sdk")
	if err != nil {
		spinner.Fail("Unable to move SDK files: ", err)
		return err
	}

	// remove sdk/sdk folder
	err = os.RemoveAll(filepath.Join("sdk", "sdk"))
	if err != nil {
		spinner.Fail("Unable to remove sdk/sdk folder: ", err)
		return err
	}

	// remove .github/ folder
	err = os.RemoveAll(filepath.Join("sdk", ".github"))
	if err != nil {
		spinner.Fail("Unable to remove .github folder: ", err)
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
