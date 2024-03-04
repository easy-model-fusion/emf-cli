package app

import (
	"github.com/easy-model-fusion/emf-cli/internal/ui"
)

const Name = "emf-cli"
const Repository = "https://github.com/easy-model-fusion"
const DownloadDirectoryPath = "./models/"

var (
	// Version is the binary version + build number
	Version string
	// BuildDate is the date of build
	BuildDate string
)

var _ui ui.UI

func Init(version, buildDate string) {
	Version = version
	BuildDate = buildDate
	initLogger()

	// Initialize the UI
	_ui = ui.NewPTermUI() // currently only pterm is supported
}

// UI returns the current UI instance
func UI() ui.UI {
	if _ui == nil {
		fatal("UI not initialized")
	}
	return _ui
}

// SetUI sets the current UI with a new one
func SetUI(newUI ui.UI) {
	_ui = newUI
}
