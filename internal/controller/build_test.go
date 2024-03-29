package controller

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"testing"
)

func TestBuildController_Build(t *testing.T) {

}

func TestBuildController_RunBuild(t *testing.T) {

}

func TestBuildController_RunBuildPyinstaller(t *testing.T) {

}

func TestBuildController_RunBuildNuitka(t *testing.T) {

}

func TestBuildController_RunBuildPyinstallerOneFile(t *testing.T) {

}

func TestBuildController_RunBuildNuitkaOneFile(t *testing.T) {

}

func TestBuildController_createModelsSymbolicLink(t *testing.T) {

}

func TestBuildController_InstallDependencies(t *testing.T) {

}

func TestCreateBuildArgs_WithEmptyCustomName(t *testing.T) {
	bc := BuildController{}
	customName := ""
	library := "pyinstaller"
	destDir := "dist"
	oneFile := false

	viper.Set("name", "testName")

	args := bc.createBuildArgs(customName, library, destDir, oneFile)
	test.AssertEqual(t, len(args), 3)
	test.AssertEqual(t, args[0], "--name=testName")
	test.AssertEqual(t, args[1], "--distpath=dist")
	test.AssertEqual(t, args[2], "main.py")
}

func TestCreateBuildArgs_WithPyInstallerAndOneFile(t *testing.T) {
	bc := BuildController{}
	customName := "custom"
	library := "pyinstaller"
	destDir := "dist"
	oneFile := true

	args := bc.createBuildArgs(customName, library, destDir, oneFile)
	test.AssertEqual(t, len(args), 4)
	test.AssertEqual(t, args[0], "-F")
	test.AssertEqual(t, args[1], "--name=custom")
	test.AssertEqual(t, args[2], "--distpath=dist")
	test.AssertEqual(t, args[3], "main.py")
}

func TestCreateBuildArgs_WithNuitkaAndOneFile(t *testing.T) {
	bc := BuildController{}
	customName := "custom"
	library := "nuitka"
	destDir := "dist"
	oneFile := true

	args := bc.createBuildArgs(customName, library, destDir, oneFile)
	test.AssertEqual(t, len(args), 5)
	test.AssertEqual(t, args[0], "-m nuitka")
	test.AssertEqual(t, args[1], "--onefile")
	test.AssertEqual(t, args[2], "--python-flag=-o custom")
	test.AssertEqual(t, args[3], "--output-dir=dist")
	test.AssertEqual(t, args[4], "main.py")

}
