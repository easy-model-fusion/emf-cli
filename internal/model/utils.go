package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"os"
	"path"
)

// ConstructConfigPaths to update the model's path to elements accordingly to its configuration.
func (m *Model) ConstructConfigPaths() {
	basePath := path.Join(app.DownloadDirectoryPath, m.Name)
	modelPath := basePath
	if m.Module == huggingface.TRANSFORMERS {
		modelPath = path.Join(modelPath, "model")
		for i, tokenizer := range m.Tokenizers {
			m.Tokenizers[i].Path = path.Join(basePath, tokenizer.Class)
		}
	}
	m.Path = modelPath
}

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

// FromHuggingfaceModel map the Huggingface API huggingface.Model to a Model
func FromHuggingfaceModel(huggingfaceModel huggingface.Model) Model {
	var model Model
	model.Name = huggingfaceModel.Name
	model.PipelineTag = huggingfaceModel.PipelineTag
	model.Module = huggingfaceModel.LibraryName
	model.Source = HUGGING_FACE
	model.Version = huggingfaceModel.LastModified
	return model
}

// BuildModelsFromDevice builds a slice of models recovered from the device folders.
func BuildModelsFromDevice() Models {

	// Get all the providers in the root folder
	providers, err := os.ReadDir(app.DownloadDirectoryPath)
	if err != nil {
		return Models{}
	}

	// Processing each provider
	var models Models
	for _, provider := range providers {
		// If it's not a directory, skip
		if !provider.IsDir() {
			continue
		}

		// Get all the models for the provider
		providerPath := path.Join(app.DownloadDirectoryPath, provider.Name())
		providerModels, err := os.ReadDir(providerPath)
		if err != nil {
			continue
		}

		// Processing each model
		for _, providerModel := range providerModels {

			// If it's not a directory, skip
			if !providerModel.IsDir() {
				continue
			}

			// Model info
			modelName := path.Join(provider.Name(), providerModel.Name())
			modelPath := path.Join(providerPath, providerModel.Name())

			// Fetching model from huggingface
			huggingfaceModel, err := app.H().GetModelById(modelName)
			if err != nil {
				// Model not found : custom
				models = append(models, Model{
					Name:            providerModel.Name(),
					Path:            modelPath,
					Source:          CUSTOM,
					AddToBinaryFile: true,
					IsDownloaded:    true,
				})
				continue
			}

			// Fetching succeeded : processing the response
			// Map API response to model.Model
			modelMapped := FromHuggingfaceModel(huggingfaceModel)

			// Leaving the version field as empty since it's impossible to trace the version back
			modelMapped.Version = ""

			// Get all the folders for the model
			directories, err := os.ReadDir(modelPath)
			if err != nil {
				continue
			}

			// Checking whether the model is empty or not
			if len(directories) == 0 {
				// Nothing to process
				continue
			}

			// Handling model by default : a special handling is required for tokenizers
			if modelMapped.Module == huggingface.DIFFUSERS {
				modelMapped.Path = modelPath
				modelMapped.AddToBinaryFile = true
				modelMapped.IsDownloaded = true
			} else if modelMapped.Module == huggingface.TRANSFORMERS {

				// Searching for the model folder and the tokenizers
				for _, directory := range directories {

					// Model folder exists : meaning the model is downloaded
					if directory.Name() == "model" {
						modelMapped.Path = path.Join(modelPath, "model")
						modelMapped.AddToBinaryFile = true
						modelMapped.IsDownloaded = true
						continue
					}

					// Otherwise : directory is considered as a tokenizer
					tokenizer := Tokenizer{
						Path:  path.Join(modelPath, directory.Name()),
						Class: directory.Name(),
					}
					modelMapped.Tokenizers = append(modelMapped.Tokenizers, tokenizer)
				}
			}

			models = append(models, modelMapped)
		}
	}

	return models
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

// Update attempts to update the model
func (m *Model) Update(mapConfigModels map[string]Model) bool {

	// Checking if the model is already configured
	_, configured := mapConfigModels[m.Name]

	// Check if model is physically present on the device
	m.ConstructConfigPaths()
	downloaded, err := m.DownloadedOnDevice()
	if err != nil {
		return false
	}

	// Process internal state of the model
	install := false
	if !configured && !downloaded {
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be added. "+
			"Would you like to add it?", m.Name))
	} else if configured && !downloaded {
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be downloaded. "+
			"Would you like to download it?", m.Name))
	} else if !configured && downloaded {
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' already exists. "+
			"Would you like to overwrite it?", m.Name))
	} else {
		// Model already configured and downloaded : a new version is available
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("New version of '%s' is available. "+
			"Would you like to overwrite its old version?", m.Name))
	}

	// Model will not be downloaded or overwritten, nothing more to do here
	if !install {
		return false
	}

	// Downloader script to skip the tokenizers download process if none selected
	var skip string

	// If transformers : select the tokenizers to update using a multiselect
	var tokenizerNames []string
	if m.Module == huggingface.TRANSFORMERS {

		// Get tokenizer names for the model
		availableNames := m.Tokenizers.GetNames()

		// Allow to select only if at least one tokenizer is available
		if len(availableNames) > 0 {

			// Prepare the tokenizers multiselect
			message := "Please select the tokenizer(s) to be updated"
			checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
			tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, availableNames, checkMark, true, true)
			app.UI().DisplaySelectedItems(tokenizerNames)

			// No tokenizer is selected : skipping so that it doesn't overwrite the default one
			if len(tokenizerNames) > 0 {
				skip = downloader.SkipValueTokenizer
			}
		}
	}

	// Prepare the script arguments
	downloaderArgs := downloader.Args{
		ModelName:         m.Name,
		ModelModule:       string(m.Module),
		ModelClass:        m.Class,
		ModelOptions:      stringutil.OptionsMapToSlice(m.Options),
		Skip:              skip,
		OnlyConfiguration: false,
	}

	// Downloading model
	success := false
	success = m.Download(downloaderArgs)
	if !success {
		// Download failed
		return false
	}

	// If transformers and at least one tokenizer were asked for an update
	if len(tokenizerNames) > 0 {

		// Bind the model tokenizers to a map for faster lookup
		mapModelTokenizers := m.Tokenizers.ToMap()

		var failedTokenizers []string
		for _, tokenizerName := range tokenizerNames {
			tokenizer := mapModelTokenizers[tokenizerName]

			// Downloading tokenizer
			success = m.DownloadTokenizer(tokenizer, downloaderArgs)
			if !success {
				// Download failed
				failedTokenizers = append(failedTokenizers, tokenizer.Class)
				continue
			}
		}

		// The process failed for at least one tokenizer
		if len(failedTokenizers) > 0 {
			pterm.Error.Println(fmt.Sprintf("The following tokenizer(s) couldn't be downloaded for '%s': %s", m.Name, failedTokenizers))
		}
	}

	return true
}

// TidyConfiguredModel downloads the missing elements that were configured
// first bool is true if success, second bool is true if model was clean from the start
func (m *Model) TidyConfiguredModel() (bool, bool) {

	// Check if model is physically present on the device
	m.ConstructConfigPaths()
	downloaded, err := m.DownloadedOnDevice()
	if err != nil {
		return false, false
	}

	// Get all the configured but not downloaded tokenizers
	missingTokenizers := m.GetTokenizersNotDownloadedOnDevice()

	// Model is clean, nothing more to do here
	if downloaded && len(missingTokenizers) == 0 {
		return true, true
	}

	// Prepare the script arguments
	downloaderArgs := downloader.Args{
		ModelName:         m.Name,
		ModelModule:       string(m.Module),
		ModelClass:        m.Class,
		ModelOptions:      stringutil.OptionsMapToSlice(m.Options),
		OnlyConfiguration: false,
	}

	// Model has yet to be downloaded
	if !downloaded {

		// Downloading model
		success := false
		success = m.Download(downloaderArgs)
		if !success {
			// Download failed
			return false, false
		}
	}

	// Some tokenizers are missing
	if len(missingTokenizers) > 0 {

		// Downloading the missing tokenizers
		var failedTokenizers []string
		for _, tokenizer := range missingTokenizers {

			// Downloading tokenizer
			success := false
			success = m.DownloadTokenizer(tokenizer, downloaderArgs)
			if !success {
				// Download failed
				failedTokenizers = append(failedTokenizers, tokenizer.Class)
				continue
			}
		}

		// The process failed for at least one tokenizer
		if len(failedTokenizers) > 0 {
			pterm.Error.Println(fmt.Sprintf("The following tokenizer(s) couldn't be downloaded for '%s': %s", m.Name, failedTokenizers))
		}
	}

	return true, false
}
