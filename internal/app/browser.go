package app

import (
	"os"
	"path"
	"strings"
)

const DownloadDirectoryPath = "./models/"

// GetDownloadedModelNames get all the downloaded models names
func GetDownloadedModelNames() (models []string, err error) {
	// Get all the models folders in the root folder
	entries, err := os.ReadDir(DownloadDirectoryPath)
	if err != nil {
		return nil, err
	}

	// Parse models
	for _, entry := range entries {
		// If it's not a directory, skip
		if !entry.IsDir() {
			continue
		}
		// Construct the model name based on directory names
		modelDirectoryPath := path.Join(DownloadDirectoryPath, entry.Name())
		err, modelFullPath := getModelFullPath(modelDirectoryPath)
		if err != nil {
			return nil, err
		}
		modelName := getModelName(path.Dir(modelDirectoryPath), modelFullPath)
		models = append(models, modelName)
	}

	return models, nil
}

// getModelFullPath gets the full path to the model's weights
func getModelFullPath(modelDirectoryPath string) (error, string) {
	// Get files and folders of the directory
	entries, err := os.ReadDir(modelDirectoryPath)
	if err != nil {
		return err, ""
	}

	// If the directory contains only a folder => folder's name is part of the model name
	if len(entries) == 1 && entries[0].IsDir() {
		// Continue the search for the model's full path
		err, foundPath := getModelFullPath(path.Join(modelDirectoryPath, entries[0].Name()))
		if err != nil {
			return err, ""
		} else {
			return nil, foundPath
		}
	} else {
		return nil, modelDirectoryPath
	}
}

// getModelName extracts the model's name from its full path
func getModelName(rootDir, fullPath string) string {
	// Remove the root directory from the full path
	name := strings.TrimPrefix(fullPath, rootDir+"/")

	return name
}
