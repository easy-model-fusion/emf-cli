package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestRunUpgrade(t *testing.T) {
	app.ReplaceUI(test.NewMockUI())
	// test "no" to the confirmation
	app.UI().(*test.MockUI).UserConfirmationResult = false

	// should not run the upgrade
	RunUpgrade([]string{})

	// test "yes" to the confirmation
	app.UI().(*test.MockUI).UserConfirmationResult = true

	// No config file, so it should return an error
	RunUpgrade([]string{})
}
