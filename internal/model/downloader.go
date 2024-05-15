package model

import (
	"context"
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"os"
	"os/signal"
	"syscall"
)

// FromDownloaderModel maps data from downloadermodel.Model to Model and keeps unchanged properties of Model.
func (m *Model) FromDownloaderModel(dlModel downloadermodel.Model) {

	// Check if ScriptModel is valid
	if !dlModel.Empty() {
		if len(dlModel.Path) != 0 {
			m.Path = fileutil.PathUniformize(dlModel.Path)
		}
		m.Module = huggingface.Module(dlModel.Module)
		m.Class = dlModel.Class
		m.Options = dlModel.Options
	}

	// Check if ScriptTokenizer is valid
	if !dlModel.Tokenizer.Empty() {

		// Mapping to tokenizer
		var tokenizer Tokenizer
		if len(dlModel.Tokenizer.Path) != 0 {
			tokenizer.Path = fileutil.PathUniformize(dlModel.Tokenizer.Path)
		}
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
func (m *Model) GetConfig(downloaderArgs downloadermodel.Args) (succeeded bool, warnings []string, err error) {
	// Add OnlyConfiguration flag to the command
	downloaderArgs.OnlyConfiguration = true

	// Running the script
	return m.executeDownload(downloaderArgs)
}

// Download attempts to download the model
func (m *Model) Download(downloaderArgs downloadermodel.Args) (succeeded bool, warnings []string, err error) {
	// Running the script
	succeeded, warnings, err = m.executeDownload(downloaderArgs)
	if err != nil {
		return false, warnings, err
	}

	if succeeded {
		m.AddToBinaryFile = !downloaderArgs.OnlyConfiguration
		m.IsDownloaded = !downloaderArgs.OnlyConfiguration
	}

	return succeeded, warnings, err
}

// DownloadTokenizer attempts to download the tokenizer
func (m *Model) DownloadTokenizer(tokenizer Tokenizer, downloaderArgs downloadermodel.Args) (success bool, warnings []string, err error) {

	// Building downloader args for the tokenizer
	downloaderArgs.SkipModel = true
	downloaderArgs.SkipTokenizer = false
	downloaderArgs.TokenizerClass = tokenizer.Class
	downloaderArgs.TokenizerOptions = stringutil.OptionsMapToSlice(tokenizer.Options)
	return m.executeDownload(downloaderArgs)

}

// executeDownload runs the download script
func (m *Model) executeDownload(downloaderArgs downloadermodel.Args) (success bool, warnings []string, err error) {
	// Running the script (with cancellation handling)
	var dlModel downloadermodel.Model

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Running the script in a goroutine (to handle cancellation, since the script can take a long time)
	go func() {
		// Running the script
		dlModel, err = app.Downloader().Execute(downloaderArgs, app.Python(), ctx)
		// Sending signal to the main goroutine that the script has finished
		done <- syscall.SIGQUIT
	}()

	switch code := <-done; {
	case code == syscall.SIGQUIT:
		// Do nothing
	case code == syscall.SIGINT:
		fallthrough
	case code == syscall.SIGTERM:
		cancel() // Cancel the context (to stop the script)
		warnings = append(warnings, "Please note that when cancelling the model, partial files may have been downloaded.")
		warnings = append(warnings, "Please remove the related model directory or the cache if you want to clean up the partial files.")
		return false, warnings, errors.New("download cancelled manually")
	}

	// make sure that the context is cancelled, even if the script has finished
	cancel()

	if err != nil {
		// Something went wrong or no data has been returned
		return false, warnings, nil
	}

	// Update the model for the configuration file
	m.FromDownloaderModel(dlModel)

	return true, warnings, err
}
