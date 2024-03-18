package app

import (
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/pterm/pterm"
	"testing"
)

func TestInitHuggingFace(t *testing.T) {
	InitHuggingFace("http://localhost:8080", "")
	test.AssertNotEqual(t, huggingFace, nil, "Should not be nil if huggingface is initialized")
	test.AssertNotEqual(t, H(), nil, "Should not be nil if huggingface is initialized")
	huggingFace = nil
	fatalCalled := false
	fatal = func(a ...interface{}) *pterm.TextPrinter {
		fatalCalled = true
		return nil
	}
	var h huggingface.HuggingFace
	test.AssertEqual(t, H(), h)
	test.AssertEqual(t, fatalCalled, true, "Should call the fatal function")
}

func TestInitGit(t *testing.T) {
	InitGit("http://localhost:8080", "")
	test.AssertNotEqual(t, gitInstance, nil, "Should not be nil if git is initialized")
	test.AssertNotEqual(t, G(), nil, "Should not be nil if git is initialized")
	gitInstance = nil
	fatalCalled := false
	fatal = func(a ...interface{}) *pterm.TextPrinter {
		fatalCalled = true
		return nil
	}
	test.AssertEqual(t, G(), nil)
	test.AssertEqual(t, fatalCalled, true, "Should call the fatal function")
}

func TestSetGit(t *testing.T) {
	// Initialize git Instance
	git := test.MockGit{Tag: "test-1.0"}

	// Set new git instance
	SetGit(&git)

	// Assertions
	test.AssertEqual(t, G(), &git)
}
