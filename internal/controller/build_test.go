package controller

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestBuildController_WithSuccessfulBuild(t *testing.T) {
	bc := BuildController{
		CustomName:     "custom",
		Library:        "pyinstaller",
		DestinationDir: "dist",
	}

	err := bc.Build("dist")
	test.AssertEqual(t, err, nil)
}

func TestBuildController_Run_WithErrors(t *testing.T) {
	app.SetUI(&mock.MockUI{})
	pythonMock := mock.MockPython{
		CalledFunctions: map[string]int{},
	}
	app.SetPython(&pythonMock)
	app.SetGit(&mock.MockGit{})

	bc := BuildController{
		CustomName:     "custom",
		Library:        "invalid",
		DestinationDir: "dist",
		OneFile:        true,
		ModelsSymlink:  false,
	}

	// config file not found
	err := bc.Run()
	test.AssertNotEqual(t, err, nil)

	// invalid library selected
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err = bc.Run()
	test.AssertEqual(t, err.Error(), "invalid library selected")

	// install dependencies error
	bc.Library = "pyinstaller"
	pythonMock.FindVEnvExecutableError = errors.New("mock error")
	err = bc.Run()
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, err.Error(), "error finding python executable: mock error")
}

func TestBuildController_Run(t *testing.T) {
	// init mocks
	app.SetUI(&mock.MockUI{})
	app.SetPython(&mock.MockPython{})
	app.SetGit(&mock.MockGit{})

	// create controller
	bc := BuildController{
		CustomName:     "custom",
		Library:        "pyinstaller",
		DestinationDir: "dist",
		OneFile:        true,
		ModelsSymlink:  false,
	}

	// create temp dir
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// run build
	err := bc.Run()
	test.AssertEqual(t, err, nil)

	// check if the dist folder has been created
	_, err = os.Stat("dist")
	test.AssertEqual(t, err, nil)

	// test with symlink creation
	bc.ModelsSymlink = true
	err = bc.Run()
	test.AssertEqual(t, err, nil)

	// check if the symlink exists
	fi, err := os.Lstat("dist/models")
	test.AssertEqual(t, err, nil)
	// check if the symlink is a symlink
	test.AssertEqual(t, fi.Mode()&os.ModeSymlink, os.ModeSymlink)
}

func TestBuildController_createModelsSymbolicLink(t *testing.T) {
	// init ui
	app.SetUI(mock.MockUI{})

	// create controller
	bc := BuildController{
		DestinationDir: "dist",
	}

	// create temp dir
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// create dist dir
	err := os.Mkdir("dist", os.ModePerm)
	test.AssertEqual(t, err, nil)

	err = bc.createModelsSymbolicLink()
	test.AssertEqual(t, err, nil)

	// check if the symlink exists
	fi, err := os.Lstat("dist/models")
	test.AssertEqual(t, err, nil)

	// check if the symlink is a symlink
	test.AssertEqual(t, fi.Mode()&os.ModeSymlink, os.ModeSymlink)
}

func TestBuildController_InstallDependencies(t *testing.T) {
	// init ui
	app.SetUI(mock.MockUI{})
	pythonMock := mock.MockPython{
		CalledFunctions: map[string]int{},
	}
	app.SetPython(pythonMock)

	// create controller
	bc := BuildController{}

	// install dependencies
	_, err := bc.InstallDependencies("pyinstaller")
	test.AssertEqual(t, err, nil)

	// check if the function has been called
	test.AssertEqual(t, pythonMock.CalledFunctions["FindVEnvExecutable"], 2)
	test.AssertEqual(t, pythonMock.CalledFunctions["ExecutePip"], 1)
}

func TestBuildController_InstallDependenciesError(t *testing.T) {
	// init ui
	app.SetUI(mock.MockUI{})
	matchErr := errors.New("mock error")
	pythonMock := mock.MockPython{
		CalledFunctions:         map[string]int{},
		FindVEnvExecutableError: matchErr,
	}
	app.SetPython(&pythonMock)

	// create controller
	bc := BuildController{}

	// install dependencies
	_, err := bc.InstallDependencies("pyinstaller")
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, err.Error(), "error finding python executable: mock error")

	// check if the function has been called
	test.AssertEqual(t, pythonMock.CalledFunctions["FindVEnvExecutable"], 1)
	test.AssertEqual(t, pythonMock.CalledFunctions["ExecutePip"], 0)

	// set execute pip error
	pythonMock.CalledFunctions = map[string]int{}
	pythonMock.FindVEnvExecutableError = nil
	pythonMock.ExecutePipError = matchErr

	// install dependencies
	_, err = bc.InstallDependencies("pyinstaller")
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, err.Error(), "error installing pyinstaller: mock error")

	// check if the function has been called
	test.AssertEqual(t, pythonMock.CalledFunctions["FindVEnvExecutable"], 2)
	test.AssertEqual(t, pythonMock.CalledFunctions["ExecutePip"], 1)
}

func TestCreateBuildArgs_WithEmptyCustomName(t *testing.T) {
	bc := BuildController{
		CustomName:     "",
		Library:        "pyinstaller",
		DestinationDir: "dist",
		OneFile:        false,
	}

	viper.Set("name", "testName")

	args := bc.createBuildArgs()
	test.AssertEqual(t, len(args), 3)
	test.AssertEqual(t, args[0], "--name=testName")
	test.AssertEqual(t, args[1], "--distpath=dist")
	test.AssertEqual(t, args[2], "main.py")
}

func TestCreateBuildArgs_WithPyInstallerAndOneFile(t *testing.T) {
	bc := BuildController{
		CustomName:     "custom",
		Library:        "pyinstaller",
		DestinationDir: "dist",
		OneFile:        true,
	}

	viper.Reset()

	args := bc.createBuildArgs()
	test.AssertEqual(t, len(args), 4)
	test.AssertEqual(t, args[0], "-F")
	test.AssertEqual(t, args[1], "--name=custom")
	test.AssertEqual(t, args[2], "--distpath=dist")
	test.AssertEqual(t, args[3], "main.py")
}

func TestCreateBuildArgs_WithNuitkaAndOneFile(t *testing.T) {
	bc := BuildController{
		CustomName:     "custom",
		Library:        "nuitka",
		DestinationDir: "dist",
		OneFile:        true,
	}

	viper.Reset()

	args := bc.createBuildArgs()
	test.AssertEqual(t, len(args), 5)
	test.AssertEqual(t, args[0], "-m nuitka")
	test.AssertEqual(t, args[1], "--onefile")
	test.AssertEqual(t, args[2], "--python-flag=-o custom")
	test.AssertEqual(t, args[3], "--output-dir=dist")
	test.AssertEqual(t, args[4], "main.py")
}
