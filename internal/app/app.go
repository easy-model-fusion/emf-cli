package app

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
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
var _python python.Python
var _downloader downloader.Downloader

func Init(version, buildDate string) {
	Version = version
	BuildDate = buildDate
	initLogger()

	// Initialize the UI
	_ui = ui.NewPTermUI() // currently only pterm is supported

	// Initialize Python
	_python = python.NewPython()

	// Initialize Downloader
	_downloader = downloader.NewScriptDownloader()
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

// Python returns the current Python instance
func Python() python.Python {
	if _python == nil {
		fatal("Python not initialized")
	}
	return _python
}

// SetPython sets the current Python with a new one
func SetPython(newPython python.Python) {
	_python = newPython
}

// Downloader returns the current Downloader instance
func Downloader() downloader.Downloader {
	return _downloader
}

// SetDownloader sets the current Downloader with a new one
func SetDownloader(newDownloader downloader.Downloader) {
	_downloader = newDownloader
}
