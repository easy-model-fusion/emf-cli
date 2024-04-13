package appselec

import (
	"github.com/easy-model-fusion/emf-cli/internal/selector"
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
var _modelSelector selector.ModelSelector

func Init(version, buildDate string) {
	Version = version
	BuildDate = buildDate

	// Initialize Selector
	_modelSelector = selector.NewTransformerModelSelector()
}

// Selector returns modelSelector instance
func Selector() selector.ModelSelector {
	return _modelSelector
}

// SetSelector sets the modelSelector instance
func SetSelector(newSelector selector.ModelSelector) {
	_modelSelector = newSelector
}
