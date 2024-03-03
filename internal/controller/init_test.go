package controller

import (
	"os"
	"testing"
)

func TestCreateProjectFiles(t *testing.T) {
	// create a temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dname)

	// change the current directory to the temporary directory
	if err = os.Chdir(dname); err != nil {
		t.Fatal(err)
	}

	// mkdir test
	err = os.Mkdir("test", os.ModePerm)
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
