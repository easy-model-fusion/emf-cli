package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

func TestRunUpgrade(t *testing.T) {
	app.SetUI(&mock.MockUI{})
	// todo: mock git?
	app.InitGit(app.Repository, "")

	// test "no" to the confirmation
	app.UI().(*mock.MockUI).UserConfirmationResult = false

	// should not run the upgrade
	RunUpgrade(false)

	// test "yes" to the confirmation
	app.UI().(*mock.MockUI).UserConfirmationResult = true

	// No config file, so it should return an error
	RunUpgrade(false)

	// create test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	config.FilePath = "."

	// upgrade should run
	RunUpgrade(false)
}
