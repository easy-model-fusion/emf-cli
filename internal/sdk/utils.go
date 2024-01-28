package sdk

import (
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// checkForUpdates Check for updates and return a whether there is an update or not
// Compares the current version of the SDK with the latest tag on GitHub
func checkForUpdates() bool {
	version := viper.GetString("sdk-tag")
	if version == "" {
		return false
	}

	tag, err := utils.GetLatestTag("sdk")
	if err != nil {
		return false
	}

	return tag != version
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
	if !canSendUpdateSuggestion() || !checkForUpdates() {
		return
	}

	pterm.DefaultBox.Println("A new version of the SDK is available!\nTo update, run 'emf-cli update'")
	pterm.Println()

	setUpdateSuggestion(true)
}
