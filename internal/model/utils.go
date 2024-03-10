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

// Empty checks if the models slice is empty.
func Empty(models []Model) bool {
	// No models currently downloaded
	return len(models) == 0
}

// ContainsByName checks if a models slice contains the requested model by name
func ContainsByName(models []Model, name string) bool {
	for _, currentModel := range models {
		if currentModel.Name == name {
			return true
		}
	}
	return false
}

// Difference returns the models in `parentSlice` that are not present in `subSlice`
func Difference(parentSlice, subSlice []Model) []Model {
	var difference []Model
	for _, item := range parentSlice {
		if !ContainsByName(subSlice, item.Name) {
			difference = append(difference, item)
		}
	}
	return difference
}

// Union returns the models present in both `slice1` and `slice2`
func Union(slice1, slice2 []Model) []Model {
	var union []Model
	for _, item := range slice1 {
		if ContainsByName(slice2, item.Name) {
			union = append(union, item)
		}
	}
	return union
}

// ModelsToMap creates a map from a slice of models for faster lookup.
func ModelsToMap(models []Model) map[string]Model {
	modelsMap := make(map[string]Model)
	for _, current := range models {
		modelsMap[current.Name] = current
	}
	return modelsMap
}

// TokenizersToMap creates a map from a slice of tokenizers for faster lookup.
func TokenizersToMap(model Model) map[string]Tokenizer {
	tokenizersMap := make(map[string]Tokenizer)
	for _, current := range model.Tokenizers {
		tokenizersMap[current.Class] = current
	}
	return tokenizersMap
}

// GetNames retrieves the names from the models.
func GetNames(models []Model) []string {
	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}
	return modelNames
}

// GetTokenizerNames retrieves the tokenizer names from a model.
func GetTokenizerNames(model Model) []string {
	var names []string
	for _, current := range model.Tokenizers {
		names = append(names, current.Class)
	}
	return names
}

// GetModelsByNames retrieves the models by their names given an input slice.
func GetModelsByNames(models []Model, namesSlice []string) []Model {
	// Create a map for faster lookup
	namesMap := stringutil.SliceToMap(namesSlice)

	// Slice of all the models that were found
	var namesModels []Model

	// Find the requested models
	for _, existingModel := range models {
		// Check if this model exists and adds it to the result
		if _, exists := namesMap[existingModel.Name]; exists {
			namesModels = append(namesModels, existingModel)
		}
	}

	return namesModels
}

// GetModelsWithSourceHuggingface return a sub-slice of models sourcing from huggingface.
func GetModelsWithSourceHuggingface(models []Model) []Model {
	var huggingfaceModels []Model
	for _, current := range models {
		if current.Source == HUGGING_FACE {
			huggingfaceModels = append(huggingfaceModels, current)
		}
	}
	return huggingfaceModels
}

// GetModelsWithIsDownloadedTrue return a sub-slice of models with isdownloaded to true.
func GetModelsWithIsDownloadedTrue(models []Model) []Model {
	var downloadedModels []Model
	for _, current := range models {
		if current.IsDownloaded {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// GetModelsWithAddToBinaryFileTrue return a sub-slice of models with AddToBinaryFile to true.
func GetModelsWithAddToBinaryFileTrue(models []Model) []Model {
	var downloadedModels []Model
	for _, current := range models {
		if current.AddToBinaryFile {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// ConstructConfigPaths to update the model's path to elements accordingly to its configuration.
func ConstructConfigPaths(current Model) Model {
	basePath := path.Join(app.DownloadDirectoryPath, current.Name)
	modelPath := basePath
	if current.Module == huggingface.TRANSFORMERS {
		modelPath = path.Join(modelPath, "model")
		for i, tokenizer := range current.Tokenizers {
			current.Tokenizers[i].Path = path.Join(basePath, tokenizer.Class)
		}
	}
	current.Path = modelPath

	return current
}

// MapToModelFromDownloaderModel maps data from downloader.Model to Model.
func MapToModelFromDownloaderModel(model Model, dlModel downloader.Model) Model {

	// Check if ScriptModel is valid
	if !downloader.EmptyModel(dlModel) {
		model.Path = stringutil.PathUniformize(dlModel.Path)
		model.Module = huggingface.Module(dlModel.Module)
		model.Class = dlModel.Class
	}

	// Check if ScriptTokenizer is valid
	if !downloader.EmptyTokenizer(dlModel.Tokenizer) {
		tokenizer := MapToTokenizerFromDownloaderTokenizer(dlModel.Tokenizer)

		// Check if tokenizer already configured and replace it
		var replaced bool
		for i := range model.Tokenizers {
			if model.Tokenizers[i].Class == tokenizer.Class {
				model.Tokenizers[i] = tokenizer
				replaced = true
			}
		}

		// Tokenizer was already found and replaced : nothing to append
		if replaced {
			return model
		}

		// Tokenizer not found : adding it to the list
		model.Tokenizers = append(model.Tokenizers, tokenizer)
	}

	return model
}

// MapToTokenizerFromDownloaderTokenizer maps data from downloader.Tokenizer to Tokenizer.
func MapToTokenizerFromDownloaderTokenizer(dlTokenizer downloader.Tokenizer) Tokenizer {
	var modelTokenizer Tokenizer
	modelTokenizer.Path = stringutil.PathUniformize(dlTokenizer.Path)
	modelTokenizer.Class = dlTokenizer.Class
	return modelTokenizer
}

// MapToModelFromHuggingfaceModel map the Huggingface API model to a model
func MapToModelFromHuggingfaceModel(huggingfaceModel huggingface.Model) Model {
	var model Model
	model.Name = huggingfaceModel.Name
	model.PipelineTag = huggingfaceModel.PipelineTag
	model.Module = huggingfaceModel.LibraryName
	model.Source = HUGGING_FACE
	model.Version = huggingfaceModel.LastModified
	return model
}

// ModelDownloadedOnDevice returns true if the model is physically present on the device.
func ModelDownloadedOnDevice(model Model) (bool, error) {

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(model.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(model.Path)
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

// TokenizerDownloadedOnDevice returns true if the tokenizer is physically present on the device.
func TokenizerDownloadedOnDevice(tokenizer Tokenizer) (bool, error) {

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(tokenizer.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(tokenizer.Path)
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

// TokenizersNotDownloadedOnDevice returns the list of tokenizers that should but are not physically present on the device.
func TokenizersNotDownloadedOnDevice(model Model) []Tokenizer {

	// Model can't have any tokenizer
	if model.Module != huggingface.TRANSFORMERS {
		return []Tokenizer{}
	}

	// Processing the configured tokenizers
	var notDownloadedTokenizers []Tokenizer
	for _, tokenizer := range model.Tokenizers {

		// Check if tokenizer is already downloaded
		downloaded, err := TokenizerDownloadedOnDevice(tokenizer)
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
func BuildModelsFromDevice() []Model {

	// Get all the providers in the root folder
	providers, err := os.ReadDir(app.DownloadDirectoryPath)
	if err != nil {
		return []Model{}
	}

	// Processing each provider
	var models []Model
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
			// TODO : class => Waiting for issue 61 to be completed : [Client] Analyze API
			modelMapped := MapToModelFromHuggingfaceModel(huggingfaceModel)

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

// Download attempts to download the model
func Download(model Model, downloaderArgs downloader.Args) (Model, bool) {
	// Running the script
	dlModel, err := downloader.Execute(downloaderArgs)

	// Something went wrong or no data has been returned
	if err != nil || dlModel.IsEmpty {
		return model, false
	}

	// Update the model for the configuration file
	model = MapToModelFromDownloaderModel(model, dlModel)
	model.AddToBinaryFile = true
	model.IsDownloaded = true

	return model, true
}

// DownloadTokenizer attempts to download the tokenizer
func DownloadTokenizer(model Model, tokenizer Tokenizer, downloaderArgs downloader.Args) (Model, bool) {

	// TODO : options tokenizer => Waiting for issue 74 to be completed : [Client] Model options to config
	// Building downloader args for the tokenizer
	downloaderArgs.Skip = downloader.SkipValueModel
	downloaderArgs.TokenizerClass = tokenizer.Class
	downloaderArgs.TokenizerOptions = []string{}

	// Running the script for the tokenizer only
	dlModelTokenizer, err := downloader.Execute(downloaderArgs)

	// Something went wrong or no data has been returned
	if err != nil || dlModelTokenizer.IsEmpty {
		return model, false
	}

	// Update the model with the tokenizer for the configuration file
	model = MapToModelFromDownloaderModel(model, dlModelTokenizer)

	return model, true
}

// Update attempts to update the model
func Update(model Model, mapConfigModels map[string]Model) bool {

	// Checking if the model is already configured
	_, configured := mapConfigModels[model.Name]

	// Check if model is physically present on the device
	model = ConstructConfigPaths(model)
	downloaded, err := ModelDownloadedOnDevice(model)
	if err != nil {
		return false
	}

	// Process internal state of the model
	install := false
	if !configured && !downloaded {
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be added. "+
			"Would you like to add it?", model.Name))
	} else if configured && !downloaded {
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be downloaded. "+
			"Would you like to download it?", model.Name))
	} else if !configured && downloaded {
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' already exists. "+
			"Would you like to overwrite it?", model.Name))
	} else {
		// Model already configured and downloaded : a new version is available
		install = app.UI().AskForUsersConfirmation(fmt.Sprintf("New version of '%s' is available. "+
			"Would you like to overwrite its old version?", model.Name))
	}

	// Model will not be downloaded or overwritten, nothing more to do here
	if !install {
		return false
	}

	// Downloader script to skip the tokenizers download process if none selected
	var skip string

	// If transformers : select the tokenizers to update using a multiselect
	var tokenizerNames []string
	if model.Module == huggingface.TRANSFORMERS {

		// Get tokenizer names for the model
		availableNames := GetTokenizerNames(model)

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

	// TODO : options model => Waiting for issue 74 to be completed : [Client] Model options to config
	// Prepare the script arguments
	downloaderArgs := downloader.Args{
		ModelName:    model.Name,
		ModelModule:  string(model.Module),
		ModelClass:   model.Class,
		ModelOptions: []string{},
		Skip:         skip,
	}

	// Downloading model
	success := false
	model, success = Download(model, downloaderArgs)
	if !success {
		// OnlyConfiguration failed
		return false
	}

	// If transformers and at least one tokenizer were asked for an update
	if len(tokenizerNames) > 0 {

		// Bind the model tokenizers to a map for faster lookup
		mapModelTokenizers := TokenizersToMap(model)

		var failedTokenizers []string
		for _, tokenizerName := range tokenizerNames {
			tokenizer := mapModelTokenizers[tokenizerName]

			// Downloading tokenizer
			success := false
			model, success = DownloadTokenizer(model, tokenizer, downloaderArgs)
			if !success {
				// OnlyConfiguration failed
				failedTokenizers = append(failedTokenizers, tokenizer.Class)
				continue
			}
		}

		// The process failed for at least one tokenizer
		if len(failedTokenizers) > 0 {
			pterm.Error.Println(fmt.Sprintf("The following tokenizer(s) couldn't be downloaded for '%s': %s", model.Name, failedTokenizers))
		}
	}

	return true
}

// TidyConfiguredModel downloads the missing elements that were configured
// first bool is true if success, second bool is true if model was clean from the start
func TidyConfiguredModel(model Model) (bool, bool) {

	// Check if model is physically present on the device
	model = ConstructConfigPaths(model)
	downloaded, err := ModelDownloadedOnDevice(model)
	if err != nil {
		return false, false
	}

	// Get all the configured but not downloaded tokenizers
	missingTokenizers := TokenizersNotDownloadedOnDevice(model)

	// Model is clean, nothing more to do here
	if downloaded && len(missingTokenizers) == 0 {
		return true, true
	}

	// TODO : options model => Waiting for issue 74 to be completed : [Client] Model options to config
	// Prepare the script arguments
	downloaderArgs := downloader.Args{
		ModelName:    model.Name,
		ModelModule:  string(model.Module),
		ModelClass:   model.Class,
		ModelOptions: []string{},
	}

	// Model has yet to be downloaded
	if !downloaded {

		// Downloading model
		success := false
		model, success = Download(model, downloaderArgs)
		if !success {
			// OnlyConfiguration failed
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
			model, success = DownloadTokenizer(model, tokenizer, downloaderArgs)
			if !success {
				// OnlyConfiguration failed
				failedTokenizers = append(failedTokenizers, tokenizer.Class)
				continue
			}
		}

		// The process failed for at least one tokenizer
		if len(failedTokenizers) > 0 {
			pterm.Error.Println(fmt.Sprintf("The following tokenizer(s) couldn't be downloaded for '%s': %s", model.Name, failedTokenizers))
		}
	}

	return true, false
}
