package controller

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"testing"
)

func TestCreateProjectFiles(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// mkdir test
	err := os.Mkdir("test", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = createProjectFiles("test", "v1.0.0")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path string
	}{
		{"test/sdk"},
		{"test/models"},
		{"test/config.yaml"},
		{"test/.gitignore"},
		{"test/main.py"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			_, err = os.Stat(test.path)
			if err != nil {
				t.Errorf("%s should exist", test.path)
			}
		})
	}

	// now test with each existing file then remove it (cover error cases)
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			err = createProjectFiles("test", "v1.0.0")
			if err == nil {
				t.Errorf("%s should return an error", test.path)
			}
			err = os.Remove(test.path)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}
