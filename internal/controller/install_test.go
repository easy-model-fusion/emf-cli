package controller

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"strings"
	"testing"
)

func TestInstallController_Run_ConfigFileMissing(t *testing.T) {
	app.SetUI(&mock.MockUI{})

	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(dname)

	ic := InstallController{}

	err = ic.Run([]string{}, false, "")
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	// FIXME: impossible to capture -> err, ok := err.(viper.ConfigFileNotFoundError)
	// so we do some nasty string contains check
	test.AssertEqual(t, strings.Contains(err.Error(), "Config File"), true, "Error should be a viper.ConfigFileNotFoundError")
}

func TestInstallController_Run_PythonNotFound(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}

	mockGit := &mock.MockGit{}
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	app.SetGit(mockGit)
	ic := InstallController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Test the case where python is not found
	err := ic.Run([]string{}, false, "")
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "python not found", "Error should be 'python not found'")
}

func TestInstallController_Run(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}
	mockPython.Success = true
	mockGit := &mock.MockGit{}
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	app.SetGit(mockGit)
	ic := InstallController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Test the case where everything is fine (cuda not used)
	err := ic.Run([]string{}, false, "")
	test.AssertEqual(t, err, nil, "Error should be nil")

	// Test the case where everything is fine (cuda used)
	err = ic.Run([]string{}, true, "")
	test.AssertEqual(t, err, nil, "Error should be nil")
}

func TestInstallController_createMissingDirectories(t *testing.T) {
	app.SetUI(&mock.MockUI{})
	ic := InstallController{}

	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Test the case where everything is fine
	err := ic.createMissingDirectories()
	test.AssertEqual(t, err, nil, "Error should be nil")

	// Test the case where sdk folder already exists
	err = ic.createMissingDirectories()
	test.AssertEqual(t, err, nil, "Error should be nil")

	// Test the case where models folder already exists
	err = ic.createMissingDirectories()
	test.AssertEqual(t, err, nil, "Error should be nil")
}

func TestInstallController_installDependencies(t *testing.T) {
	mockPython := &mock.MockPython{
		CalledFunctions: make(map[string]int),
	}
	mockPython.Success = true
	app.SetUI(&mock.MockUI{})
	app.SetPython(mockPython)
	ic := InstallController{}

	// Test the case where everything is fine
	err := ic.installDependencies("python", false)
	test.AssertEqual(t, err, nil, "Error should be nil")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 2, "FindVEnvExecutable should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 0, "CreateVirtualEnv should not be called")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 2, "InstallDependencies should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 0, "ExecutePip should not be called")

	mockPython.CalledFunctions = make(map[string]int)

	// Test the case where venv is already installed
	err = ic.installDependencies("python", false)
	test.AssertEqual(t, err, nil, "Error should be nil")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 2, "FindVEnvExecutable should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 0, "CreateVirtualEnv should not be called")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 2, "InstallDependencies should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 0, "ExecutePip should not be called")

	// Test the case where findvenvexec has error
	mockPython.FindVEnvExecutableError = errors.New("pip not found")
	mockPython.CalledFunctions = make(map[string]int)

	err = ic.installDependencies("python", false)
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "pip not found", "Error should be 'pip not found'")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 2, "FindVEnvExecutable should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 1, "CreateVirtualEnv should be called")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 0, "InstallDependencies should not be called")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 0, "ExecutePip should not be called")

	// Test the case where findvenvexec has error & createvirtualenv has also error
	mockPython.CreateVirtualEnvError = errors.New("error creating virtual env")
	mockPython.CalledFunctions = make(map[string]int)

	err = ic.installDependencies("python", false)
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "error creating virtual env", "Error should be 'error creating virtual env'")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 1, "FindVEnvExecutable should be called once")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 1, "CreateVirtualEnv should be called once")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 0, "InstallDependencies should not be called")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 0, "ExecutePip should not be called")

	mockPython.FindVEnvExecutableError = nil
	mockPython.CreateVirtualEnvError = nil
	mockPython.CalledFunctions = make(map[string]int)

	// Test the case where install dependencies has error
	mockPython.InstallDependenciesError = errors.New("error installing dependencies")

	err = ic.installDependencies("python", false)
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "error installing dependencies", "Error should be 'error installing dependencies'")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 2, "FindVEnvExecutable should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 0, "CreateVirtualEnv should not be called")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "InstallDependencies should be called once")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 0, "ExecutePip should not be called")

	// Test the case where executepip has error
	mockPython.ExecutePipError = errors.New("error executing pip")
	mockPython.InstallDependenciesError = nil
	mockPython.CalledFunctions = make(map[string]int)

	err = ic.installDependencies("python", true)
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "error executing pip", "Error should be 'error executing pip'")
	test.AssertEqual(t, mockPython.CalledFunctions["FindVEnvExecutable"], 2, "FindVEnvExecutable should be called twice")
	test.AssertEqual(t, mockPython.CalledFunctions["CreateVirtualEnv"], 0, "CreateVirtualEnv should not be called")
	test.AssertEqual(t, mockPython.CalledFunctions["InstallDependencies"], 1, "InstallDependencies should be called once")
	test.AssertEqual(t, mockPython.CalledFunctions["ExecutePip"], 1, "ExecutePip should be called once")
}

func TestInstallController_cloneSDK(t *testing.T) {
	mockGit := &mock.MockGit{}
	app.SetUI(&mock.MockUI{})
	app.SetGit(mockGit)
	ic := InstallController{}

	// Test the case where everything is fine
	err := ic.cloneSDK()
	test.AssertEqual(t, err, nil, "Error should be nil")

	// Test the case where git clone has error
	mockGit.CloneSDKError = errors.New("error cloning sdk")

	err = ic.cloneSDK()
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "error cloning sdk", "Error should be 'error cloning sdk'")
}

// TestInstallController_cloneSDKConfirm tests the cloneSDKConfirm function.
func TestInstallController_cloneSDKConfirm_retry(t *testing.T) {
	// Init
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	mockGit := &mock.MockGit{}
	ui := &mock.MockUI{}
	app.SetUI(&mock.MockUI{})
	app.SetGit(mockGit)
	app.SetUI(ui)
	ic := InstallController{}

	// Test the case where git clone has error and retry is true
	mockGit.CloneSDKError = errors.New("error cloning sdk")

	ui.UserConfirmationResult = true

	err := ic.cloneSDK()
	test.AssertNotEqual(t, err, nil, "Error should not be nil")
	test.AssertEqual(t, err.Error(), "error cloning sdk", "Error should be 'error cloning sdk'")
}
