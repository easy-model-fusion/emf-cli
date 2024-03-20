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
	ic := InitController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// mkdir test
	err := os.Mkdir("test", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ic.createProjectFiles("test", "v1.0.0")
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
			err = ic.createProjectFiles("test", "v1.0.0")
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
	ic := InitController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := ic.createProjectFolder("test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("test")
	if err != nil {
		t.Errorf("test should exist")
	}

	err = ic.createProjectFolder("test")
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
	ic := InitController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// first test with no test/ directory created
	err := ic.cloneSDK("test", "")
	test.AssertNotEqual(t, err, nil, "Expected error")

	// create the test directory
	err = os.Mkdir("test", os.ModePerm)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = ic.cloneSDK("test", "")
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
	err = ic.cloneSDK("test", "v1.0.0")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	recreateProjectFolder(t)

	// test with getlatesttag error
	mockGit.LatestTagError = errors.New("LatestTagError")

	err = ic.cloneSDK("test", "")
	t.Logf("%v", err)
	test.AssertNotEqual(t, err, nil, "Expected error")

	mockGit.LatestTagError = nil
	mockGit.CloneSDKError = errors.New("CloneSDKError")

	err = ic.cloneSDK("test", "")
	t.Logf("%v", err)
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "CloneSDKError", "Expected CloneSDKError error")
}

func TestInstallDependencies(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	ic := InitController{}

	// first test should pass without any errors (torch cuda not used)
	err := ic.installDependencies("test", false)
	test.AssertEqual(t, err, nil, "Expected no error")
	// torch cuda not used, we should have only one call to InstallDependencies, and one to FindVEnvExecutable
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "Expected 1 call to InstallDependencies")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 1, "Expected 1 call to FindVEnvExecutable")

	// reset called functions
	mockPython.CalledFunctions = make(map[string]int)

	// second test should pass without any errors (torch cuda used)
	err = ic.installDependencies("test", true)
	test.AssertEqual(t, err, nil, "Expected no error")
	// torch cuda used, we should have one call to InstallDependencies, one to FindVEnvExecutable, and two to ExecutePip
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "Expected 1 call to InstallDependencies")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 1, "Expected 1 call to FindVEnvExecutable")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 2, "Expected 2 calls to ExecutePip")
}

func TestInstallDependencies_WithErrors(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	ic := InitController{}

	// test with error finding venv executable (torch cuda not used)
	mockPython.FindVEnvExecutableError = errors.New("FindVEnvExecutableError")

	err := ic.installDependencies("test", false)
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "FindVEnvExecutableError")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 0, "Expected 0 calls to InstallDependencies")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 1, "Expected 1 call to FindVEnvExecutable")

	mockPython.CalledFunctions = make(map[string]int)
	// test with error installing dependencies (torch cuda not used)
	mockPython.FindVEnvExecutableError = nil
	mockPython.InstallDependenciesError = errors.New("InstallDependenciesError")

	err = ic.installDependencies("test", false)
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "InstallDependenciesError")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "Expected 1 call to InstallDependencies")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 1, "Expected 1 call to FindVEnvExecutable")

	mockPython.CalledFunctions = make(map[string]int)
	mockPython.InstallDependenciesError = nil
	// test with error executing pip (torch cuda used)
	mockPython.ExecutePipError = errors.New("ExecutePipError")

	err = ic.installDependencies("test", true)
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "ExecutePipError")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "Expected 1 call to InstallDependencies")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 1, "Expected 1 call to FindVEnvExecutable")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 1, "Expected 1 call to ExecutePip")
}

func TestCreateProject(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}
	mockGit := &mock.MockGit{}
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	app.SetGit(mockGit)
	ic := InitController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// create test directory
	err := os.Mkdir("test", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// test with no test/ directory created
	err = ic.createProject("test", false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")

	err = os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}

	// test with no python
	mockPython.Success = false

	err = ic.createProject("test", false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "python not found")

	err = os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}

	mockPython.Success = true

	// test
	err = ic.createProject("test", true, "")
	test.AssertEqual(t, err, nil, "Expected no error")

	// check if files are created
	testFilesCreated(t)

	err = os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}

	// test with create virtual env error
	mockPython.CalledFunctions = make(map[string]int)
	mockPython.CreateVirtualEnvError = errors.New("CreateVirtualEnvError")

	err = ic.createProject("test", false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "CreateVirtualEnvError")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 1, "Expected 1 call to CreateVirtualEnv")

	err = os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}

	// test with clone sdk error
	mockGit.CloneSDKError = errors.New("CloneSDKError")

	err = ic.createProject("test", false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "CloneSDKError", "Expected CloneSDKError error")

	err = os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}

	// test with install dependencies error
	mockGit.CloneSDKError = nil
	mockPython.CalledFunctions = make(map[string]int)
	mockPython.CreateVirtualEnvError = nil
	mockPython.InstallDependenciesError = errors.New("InstallDependenciesError")

	err = ic.createProject("test", false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "InstallDependenciesError", "Expected InstallDependenciesError error")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "Expected 1 call to InstallDependencies")
}

func TestInitController_Run(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}
	mockGit := &mock.MockGit{}
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	app.SetGit(mockGit)
	ic := InitController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	app.UI().(*mock.MockUI).UserInputResult = "hey"

	// test with create project error (python not found)
	err := ic.Run([]string{}, false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")
	test.AssertEqual(t, err.Error(), "python not found")

	// folder should have been removed
	_, err = os.Stat("hey")
	test.AssertNotEqual(t, err, nil, "Expected error")

	// test with create project with args name & folder already exists
	err = os.Mkdir("hey", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run([]string{"hey"}, false, "")
	test.AssertNotEqual(t, err, nil, "Expected error")

	// folder should not have been removed
	_, err = os.Stat("hey")
	test.AssertEqual(t, err, nil, "Expected no error")

	mockPython.Success = true

	// test success
	err = ic.Run([]string{"test"}, false, "")
	test.AssertEqual(t, err, nil, "Expected no error")

	// check if files are created
	testFilesCreated(t)
}

func testFilesCreated(t *testing.T) {
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
			_, err := os.Stat(testInstance.path)
			if err != nil {
				t.Errorf("%s should exist", testInstance.path)
			}
		})
	}
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
