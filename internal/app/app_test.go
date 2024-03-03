package app

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestInit(t *testing.T) {
	Init("1.0.0", "2021-01-01")
	if Version != "1.0.0" {
		t.Errorf("Version should be 1.0.0, got %s", Version)
	}
	if BuildDate != "2021-01-01" {
		t.Errorf("BuildDate should be 2021-01-01, got %s", BuildDate)
	}

	test.AssertNotEqual(t, _ui, nil, "UI should not be nil")
}

func TestUI(t *testing.T) {
	Init("1.0.0", "2021-01-01")

	test.AssertNotEqual(t, UI(), nil, "UI should not be nil")

	SetUI(nil)

	test.AssertEqual(t, UI(), nil, "UI should be nil")
}
