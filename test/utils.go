package test

import (
	"github.com/easy-model-fusion/client/sdk"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

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
	//err = os.Chdir(dname)
	//if err != nil {
	//	t.Error(err)
	//}

	// Create config file
	file, err := os.Create("config.yaml")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = file.Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	content, err := fs.ReadFile(sdk.EmbeddedFiles, "config.yaml")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = os.WriteFile(filepath.Join(dname, "config.yaml"), content, os.ModePerm)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return dname
}
