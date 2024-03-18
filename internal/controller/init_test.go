package controller

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	mock "github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"
)

func TestCreateProjectFiles(t *testing.T) {
	app.SetUI(&mock.MockUI{})
	app.SetGit(&mock.MockGit{})

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
	for _, testInstance := range tests {
		t.Run(testInstance.path, func(t *testing.T) {
			_, err = os.Stat(testInstance.path)
			if err != nil {
				t.Errorf("%s should exist", testInstance.path)
			}
		})
	}

	// now test with each existing file then remove it (cover error cases)
	for _, testInstance := range tests {
		t.Run(testInstance.path, func(t *testing.T) {
			err = createProjectFiles("test", "v1.0.0")
			if err == nil {
				t.Errorf("%s should return an error", testInstance.path)
			}
			err = os.Remove(testInstance.path)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}

func TestCreateProjectFolder(t *testing.T) {
	app.SetUI(&mock.MockUI{})
	app.SetGit(&mock.MockGit{})

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := createProjectFolder("test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("test")
	if err != nil {
		t.Errorf("test should exist")
	}

	err = createProjectFolder("test")
	if err == nil {
		t.Errorf("test should return an error")
	}
	err = os.Remove("test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCloneSDK(t *testing.T) {
	mockGit := &mock.MockGit{}
	app.SetUI(&mock.MockUI{})
	app.SetGit(mockGit)

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// first test with no test/ directory created
	err := cloneSDK("test", "")
	test.AssertNotEqual(t, err, nil, "Expected error")

	// create the test directory
	err = os.Mkdir("test", os.ModePerm)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = cloneSDK("test", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = os.Stat("test/sdk/.git")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	recreateProjectFolder(t)

	// test with custom tag
	err = cloneSDK("test", "v1.0.0")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	recreateProjectFolder(t)

	// test with getlatesttag error
	mockGit.LatestTagError = errors.New("LatestTagError")

	err = cloneSDK("test", "")
	t.Logf("%v", err)
	test.AssertNotEqual(t, err, nil, "Expected error")

	mockGit.LatestTagError = nil
	mockGit.CloneSDKError = errors.New("CloneSDKError")

	err = cloneSDK("test", "")
	t.Logf("%v", err)
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "CloneSDKError", "Expected CloneSDKError error")
}

// recreateProjectFolder removes the test/ directory and creates it again
func recreateProjectFolder(t *testing.T) {
	err := os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir("test", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
