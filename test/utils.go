package test

import (
	"github.com/easy-model-fusion/emf-cli/sdk"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// checkErrDeleteFolder Check if an error is not nil, delete the folder and fail the test
func checkErrDeleteFolder(t *testing.T, err error, dname string) {
	if err == nil {
		return
	}
	t.Error(err)
	err = os.RemoveAll(dname)
	if err != nil {
		t.Error(err)
	}
	t.FailNow()
}

// CreateFullTestSuite Create a full test suite
// Please delete the directory after use (defer os.RemoveAll(dname))
func CreateFullTestSuite(t *testing.T) (directoryPath string) {
	// Create temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Chdir to a temporary directory
	err = os.Chdir(dname)
	checkErrDeleteFolder(t, err, dname)

	// Create config file from embedded file
	content, err := fs.ReadFile(sdk.EmbeddedFiles, "config.yaml")
	checkErrDeleteFolder(t, err, dname)

	err = os.WriteFile("config.yaml", content, os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	return dname
}

// CreateFullTestSuite Create a full test suite
// Please delete the directory after use (defer os.RemoveAll(dname))
func CreateModelsFolderFullTestSuite(t *testing.T) (directoryPath string) {
	// Create temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Chdir to a temporary directory
	err = os.Chdir(dname)
	checkErrDeleteFolder(t, err, dname)

	//create models repository
	err = os.Mkdir("models", os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	// Create mock models
	err = os.MkdirAll(filepath.Join("models", "model1", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	err = os.MkdirAll(filepath.Join("models", "model2", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	err = os.MkdirAll(filepath.Join("models", "model3"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	return dname
}
