package app

import (
	"github.com/easy-model-fusion/emf-cli/internal/ui"
)

const Name = "emf-cli"
const Repository = "https://github.com/easy-model-fusion"

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
	return _ui
}

// ReplaceUI replaces the current UI with a new one
func ReplaceUI(newUI ui.UI) {
	_ui = newUI
}
