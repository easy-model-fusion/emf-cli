package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.Init("", "")
	os.Exit(m.Run())
}

func TestRunInit(t *testing.T) {
}
