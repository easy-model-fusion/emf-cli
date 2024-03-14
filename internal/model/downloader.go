package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

// FromDownloaderModel maps data from downloader.Model to Model and keeps unchanged properties of Model.
func (m *Model) FromDownloaderModel(dlModel downloader.Model) {

	// Check if ScriptModel is valid
	if !downloader.EmptyModel(dlModel) {
		m.Path = stringutil.PathUniformize(dlModel.Path)
		m.Module = huggingface.Module(dlModel.Module)
		m.Class = dlModel.Class
		m.Options = dlModel.Options
	}

	// Check if ScriptTokenizer is valid
	if !downloader.EmptyTokenizer(dlModel.Tokenizer) {

		// Mapping to tokenizer
		var tokenizer Tokenizer
		tokenizer.Path = stringutil.PathUniformize(dlModel.Tokenizer.Path)
		tokenizer.Class = dlModel.Tokenizer.Class
		tokenizer.Options = dlModel.Tokenizer.Options

		// Check if tokenizer already configured and replace it
		var replaced bool
		for i := range m.Tokenizers {
			if m.Tokenizers[i].Class == tokenizer.Class {
				m.Tokenizers[i] = tokenizer
				replaced = true
			}
		}

		// Tokenizer was already found and replaced : nothing to append
		if replaced {
			return
		}

		// Tokenizer not found : adding it to the list
		m.Tokenizers = append(m.Tokenizers, tokenizer)
	}
}

// GetConfig attempts to get the model's configuration
func (m *Model) GetConfig(downloaderArgs downloader.Args) bool {
	// Add OnlyConfiguration flag to the command
	downloaderArgs.OnlyConfiguration = true

	// Running the script
	dlModel, err := downloader.Execute(downloaderArgs)

	// Something went wrong or no data has been returned
	if err != nil || dlModel.IsEmpty {
		return false
	}

	// Update the model for the configuration file
	m.FromDownloaderModel(dlModel)

	return true
}

// Download attempts to download the model
func (m *Model) Download(downloaderArgs downloader.Args) bool {
	// Running the script
	dlModel, err := downloader.Execute(downloaderArgs)

	// Something went wrong or no data has been returned
	if err != nil || dlModel.IsEmpty {
		return false
	}

	// Update the model for the configuration file
	m.FromDownloaderModel(dlModel)
	m.AddToBinaryFile = !downloaderArgs.OnlyConfiguration
	m.IsDownloaded = !downloaderArgs.OnlyConfiguration

	return true
}

// DownloadTokenizer attempts to download the tokenizer
func (m *Model) DownloadTokenizer(tokenizer Tokenizer, downloaderArgs downloader.Args) bool {

	// Building downloader args for the tokenizer
	downloaderArgs.Skip = downloader.SkipValueModel
	downloaderArgs.TokenizerClass = tokenizer.Class
	downloaderArgs.TokenizerOptions = stringutil.OptionsMapToSlice(tokenizer.Options)

	// Running the script for the tokenizer only
	dlModel, err := downloader.Execute(downloaderArgs)

	// Something went wrong or no data has been returned
	if err != nil || dlModel.IsEmpty {
		return false
	}

	// Update the model with the tokenizer for the configuration file
	m.FromDownloaderModel(dlModel)

	return true
}
