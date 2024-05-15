package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/dotenv"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"path/filepath"
	"strings"
)

type Models []Model
type Model struct {
	Name            string
	Path            string
	Module          huggingface.Module
	Class           string
	Options         map[string]string
	Tokenizers      Tokenizers
	PipelineTag     huggingface.PipelineTag
	Source          string
	AddToBinaryFile bool
	IsDownloaded    bool
	Version         string
	AccessToken     string
}

type Tokenizers []Tokenizer
type Tokenizer struct {
	Path    string
	Class   string
	Options map[string]string
}

// Sources
const (
	HUGGING_FACE = "hugging_face"
	CUSTOM       = "custom"
)

// Empty checks if the models slice is empty.
func (m Models) Empty() bool {
	return len(m) == 0
}

// ContainsByName checks if a models slice contains the requested model name
func (m Models) ContainsByName(name string) bool {
	for _, currentModel := range m {
		if currentModel.Name == name {
			return true
		}
	}
	return false
}

// Difference returns the models in that are not present in `slice`
func (m Models) Difference(slice Models) Models {
	var difference Models
	for _, item := range m {
		if !slice.ContainsByName(item.Name) {
			difference = append(difference, item)
		}
	}
	return difference
}

// Union returns the models present in `slice` as well
func (m Models) Union(slice Models) Models {
	var union Models
	for _, item := range m {
		if slice.ContainsByName(item.Name) {
			union = append(union, item)
		}
	}
	return union
}

// Map creates a map from models for faster lookup.
func (m Models) Map() map[string]Model {
	modelsMap := make(map[string]Model)
	for _, current := range m {
		modelsMap[current.Name] = current
	}
	return modelsMap
}

// Map creates a map from tokenizers for faster lookup.
func (t Tokenizers) Map() map[string]Tokenizer {
	tokenizersMap := make(map[string]Tokenizer)
	for _, current := range t {
		tokenizersMap[current.Class] = current
	}
	return tokenizersMap
}

// GetNames retrieves the names from the models.
func (m Models) GetNames() []string {
	var modelNames []string
	for _, item := range m {
		modelNames = append(modelNames, item.Name)
	}
	return modelNames
}

// GetNames retrieves the names from the tokenizers.
func (t Tokenizers) GetNames() []string {
	var tokenizerNames []string
	for _, current := range t {
		tokenizerNames = append(tokenizerNames, current.Class)
	}
	return tokenizerNames
}

// FilterWithNames retrieves the models by their names given an input slice.
func (m Models) FilterWithNames(namesSlice []string) Models {
	// Create a map for faster lookup
	namesMap := stringutil.SliceToMap(namesSlice)

	// Slice of all the models that were found
	var namesModels Models

	// Find the requested models
	for _, existingModel := range m {
		// Check if this model exists and adds it to the result
		if _, exists := namesMap[existingModel.Name]; exists {
			namesModels = append(namesModels, existingModel)
		}
	}

	return namesModels
}

// FilterWithClass retrieves the tokenizers by their class given an input slice.
func (t Tokenizers) FilterWithClass(namesSlice []string) Tokenizers {
	// Create a map for faster lookup
	namesMap := stringutil.SliceToMap(namesSlice)

	// Slice of all the Tokenizers that were found
	var namesTokenizers Tokenizers

	// Find the requested Tokenizer
	for _, existingTokenizer := range t {
		// Check if this tokenizer exists and adds it to the result
		if _, exists := namesMap[existingTokenizer.Class]; exists {
			namesTokenizers = append(namesTokenizers, existingTokenizer)
		}
	}
	return namesTokenizers
}

// FilterWithSourceHuggingface return a sub-slice of models sourcing from huggingface.
func (m Models) FilterWithSourceHuggingface() Models {
	var huggingfaceModels Models
	for _, current := range m {
		if current.Source == HUGGING_FACE {
			huggingfaceModels = append(huggingfaceModels, current)
		}
	}
	return huggingfaceModels
}

// FilterWithIsDownloadedOrAddToBinaryFileTrue return a sub-slice of models with IsDownloaded or AddToBinaryFile  to true.
func (m Models) FilterWithIsDownloadedOrAddToBinaryFileTrue() Models {
	var downloadedModels Models
	for _, current := range m {
		if current.IsDownloaded || current.AddToBinaryFile {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// FilterWithIsDownloadedTrue return a sub-slice of models with IsDownloaded to true.
func (m Models) FilterWithIsDownloadedTrue() Models {
	var downloadedModels Models
	for _, current := range m {
		if current.IsDownloaded {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// FilterWithAddToBinaryFileTrue return a sub-slice of models with AddToBinaryFile to true.
func (m Models) FilterWithAddToBinaryFileTrue() Models {
	var downloadedModels Models
	for _, current := range m {
		if current.AddToBinaryFile {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// GetBasePath return the base path to the model
func (m *Model) GetBasePath() string {
	return filepath.Join(app.DownloadDirectoryPath, m.Name)
}

// UpdatePaths to update the model's path to elements accordingly to its configuration.
func (m *Model) UpdatePaths() {
	if m.Path == "" {
		basePath := m.GetBasePath()
		modelPath := basePath
		if m.Module == huggingface.TRANSFORMERS {
			modelPath = filepath.Join(modelPath, "model")
			for i, tokenizer := range m.Tokenizers {
				m.Tokenizers[i].Path = filepath.Join(basePath, tokenizer.Class)
			}
		}
		m.Path = modelPath
	}
}

// Difference returns the models in that are not present in `slice`
func (t Tokenizers) Difference(slice Tokenizers) Tokenizers {
	var difference Tokenizers
	for _, item := range t {
		if !slice.ContainsByClass(item.Class) {
			difference = append(difference, item)
		}
	}
	return difference
}

// ContainsByClass checks if a tokenizers slice contains the requested tokenizers name
func (t Tokenizers) ContainsByClass(class string) bool {
	for _, currentTokenizer := range t {
		if currentTokenizer.Class == class {
			return true
		}
	}
	return false
}

// SetAccessTokenKey sets unique access token key for model
func (m *Model) setAccessTokenKey() error {
	// Convert model name to upper case
	key := strings.ToUpper(m.Name)
	// Replace "/", "." and "-" with "_"
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, "-", "_")

	// Prepend "ACCESS_TOKEN_" to the key
	key = "ACCESS_TOKEN_" + key

	// Check for duplicates
	key, err := dotenv.SetNewEnvKey(key)
	if err != nil {
		return err
	}
	m.AccessToken = key
	return nil
}

// SaveAccessToken saves the access token in the .env file
func (m *Model) SaveAccessToken(accessToken string) error {
	err := m.setAccessTokenKey()
	if err != nil {
		return err
	}
	return dotenv.AddNewEnvVariable(m.AccessToken, accessToken)
}

// GetAccessToken gets the access token from the .env file
func (m *Model) GetAccessToken() (string, error) {
	return dotenv.GetEnvValue(m.AccessToken)
}
