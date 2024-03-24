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

type TestSuite struct {
	dname   string // temporary directory name (and new working directory during the test)
	oldWd   string // old working directory (to go back to it after the test)
	created bool   // whether the temporary directory has been created
}

// CreateFullTestSuite Create a full test suite
// Please clean test suite after use (defer ts.CleanTestSuite())
func (ts *TestSuite) CreateFullTestSuite(t *testing.T) (directoryPath string) {
	// Create temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ts.created = true
	ts.dname = dname

	// Save the old working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ts.oldWd = oldWd

	// Chdir to a temporary directory
	err = os.Chdir(dname)
	checkErrDeleteFolder(t, err, dname)

	// Create config file from embedded file
	content, err := fs.ReadFile(sdk.EmbeddedFiles, "config.yaml")
	checkErrDeleteFolder(t, err, dname)

	err = os.WriteFile("config.yaml", content, os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	// Create generated code file from embedded file
	err = os.Mkdir("sdk", os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	genFile, err := os.Create("sdk/generated_models.py")
	checkErrDeleteFolder(t, err, dname)
	err = genFile.Close()
	checkErrDeleteFolder(t, err, dname)
	genFile, err = os.Create(".env")
	checkErrDeleteFolder(t, err, dname)
	err = genFile.Close()
	checkErrDeleteFolder(t, err, dname)

	return dname
}

const FullTestSuiteModelsCount = 4

// CreateModelsFolderFullTestSuite Create a full test suite
// Please delete the directory after use (defer test.Clean())
func (ts *TestSuite) CreateModelsFolderFullTestSuite(t *testing.T) (directoryPath string) {
	// Create temporary directory
	dname := ts.CreateFullTestSuite(t)

	//create models repository
	err := os.Mkdir("models", os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	// Create mock models (if you change this, change FullTestSuiteModelsCount)
	err = os.MkdirAll(filepath.Join("models", "model1", "name", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	err = os.MkdirAll(filepath.Join("models", "model2", "name", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	err = os.MkdirAll(filepath.Join("models", "model3", "name", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	err = os.MkdirAll(filepath.Join("models", "model4", "name", "model", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)
	err = os.MkdirAll(filepath.Join("models", "model4", "name", "tokenizer", "weights"), os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	return dname
}

// CleanTestSuite Clean a test suite
func (ts *TestSuite) CleanTestSuite(t *testing.T) {
	// go back to the old working directory
	err := os.Chdir(ts.oldWd)
	if err != nil {
		t.Error(err)
	}

	// remove the temporary directory
	err = os.RemoveAll(ts.dname)
	if err != nil {
		t.Error(err)
	}
}
