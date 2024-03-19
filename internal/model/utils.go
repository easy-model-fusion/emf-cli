package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"os"
	"path"
)

// DownloadedOnDevice returns true if the model is physically present on the device.
func (m *Model) DownloadedOnDevice(useBasePath bool) (bool, error) {

	// Adapt the model path
	modelPath := m.Path
	if useBasePath {
		modelPath = m.GetBasePath()
	}

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(modelPath)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(modelPath)
	if err != nil {
		// An error occurred
		return false, err
	} else if empty {
		// Model is not downloaded on the device
		return false, nil
	}

	// Model is downloaded on the device
	return true, nil
}

// DownloadedOnDevice returns true if the tokenizer is physically present on the device.
func (t *Tokenizer) DownloadedOnDevice() (bool, error) {

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(t.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(t.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if empty {
		// Model is not downloaded on the device
		return false, nil
	}

	// Model is downloaded on the device
	return true, nil
}

// GetTokenizersNotDownloadedOnDevice returns the list of tokenizers that should but are not physically present on the device.
func (m *Model) GetTokenizersNotDownloadedOnDevice() Tokenizers {

	// Model can't have any tokenizer
	if m.Module != huggingface.TRANSFORMERS {
		return Tokenizers{}
	}

	// Processing the configured tokenizers
	var notDownloadedTokenizers Tokenizers
	for _, tokenizer := range m.Tokenizers {

		// Check if tokenizer is already downloaded
		downloaded, err := tokenizer.DownloadedOnDevice()
		if err != nil {
			// An error occurred
			continue
		} else if !downloaded {
			// Tokenizer is not downloaded on the device
			notDownloadedTokenizers = append(notDownloadedTokenizers, tokenizer)
			continue
		}
	}

	return notDownloadedTokenizers
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

// Update attempts to update the model
func (m *Model) Update(mapConfigModels map[string]Model) bool {

	// Checking if the model is already configured
	_, configured := mapConfigModels[m.Name]

	// Check if model is physically present on the device
	m.UpdatePaths()
	downloaded, err := m.DownloadedOnDevice(false)
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
		mapModelTokenizers := m.Tokenizers.Map()

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

// UpdateTokenizer attempts to update the tokenizers.
func (m *Model) UpdateTokenizer(
	selectedTokenizerNames []string,
) bool {

	// Check if model is physically present on the device
	m.UpdatePaths()
	downloaded, err := m.DownloadedOnDevice(false)
	if err != nil {
		return false
	}

	// Process internal state of the model
	install := false
	if !downloaded {
		print("Model '%s' has yet to be " +
			"added or downloaded. ")
		return false
	}

	install = app.UI().AskForUsersConfirmation(m.Name)

	// Model will not be downloaded or overwritten, nothing more to do here
	if !install {
		return false
	}

	// Downloader script to skip the tokenizers download process if none selected
	var skip string

	// Prepare the script arguments
	downloaderArgs := downloader.Args{
		ModelName:         m.Name,
		ModelModule:       string(m.Module),
		ModelClass:        m.Class,
		ModelOptions:      stringutil.OptionsMapToSlice(m.Options),
		Skip:              skip,
		OnlyConfiguration: false,
	}

	success := false

	// If transformers and at least one tokenizer were asked for an update
	if len(selectedTokenizerNames) > 0 {

		// Bind the model tokenizers to a map for faster lookup
		mapModelTokenizers := m.Tokenizers.Map()

		var failedTokenizers []string
		for _, tokenizerName := range selectedTokenizerNames {
			tokenizer := mapModelTokenizers[tokenizerName]

			// Downloading tokenizer
			print("downloading tokenizer")

			success = m.DownloadTokenizer(tokenizer, downloaderArgs)
			if !success {
				print("download failed")
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
	m.UpdatePaths()
	downloaded, err := m.DownloadedOnDevice(false)
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
