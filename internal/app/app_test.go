package app

import (
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/pterm/pterm"
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

	fatalCalled := false
	fatal = func(a ...interface{}) *pterm.TextPrinter {
		fatalCalled = true
		return nil
	}
	var u ui.UI
	test.AssertEqual(t, UI(), u)
	test.AssertEqual(t, fatalCalled, true, "Should call the fatal function")
}

func TestPython(t *testing.T) {
	Init("1.0.0", "2021-01-01")

	test.AssertNotEqual(t, Python(), nil, "Python should not be nil")

	SetPython(nil)

	fatalCalled := false
	fatal = func(a ...interface{}) *pterm.TextPrinter {
		fatalCalled = true
		return nil
	}
	var py python.Python
	test.AssertEqual(t, Python(), py)
	test.AssertEqual(t, fatalCalled, true, "Should call the fatal function")
}
