package app

import (
	git2 "github.com/easy-model-fusion/emf-cli/internal/git"
	"github.com/easy-model-fusion/emf-cli/internal/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/pterm/pterm"
	"testing"
)

func TestInitHuggingFace(t *testing.T) {
	InitHuggingFace("http://localhost:8080", "")
	test.AssertNotEqual(t, huggingFace, nil, "Should not be nil if huggingface is initialized")
	test.AssertNotEqual(t, H(), nil, "Should not be nil if huggingface is initialized")
	test.AssertEqual(t, H().BaseUrl, "http://localhost:8080", "Should be equal to the base url")
	huggingFace = nil
	fatalCalled := false
	fatal = func(a ...interface{}) *pterm.TextPrinter {
		fatalCalled = true
		return nil
	}
	var h *huggingface.HuggingFace
	test.AssertEqual(t, H(), h)
	test.AssertEqual(t, fatalCalled, true, "Should call the fatal function")
}

func TestInitGit(t *testing.T) {
	InitGit("http://localhost:8080", "")
	test.AssertNotEqual(t, git, nil, "Should not be nil if git is initialized")
	test.AssertNotEqual(t, G(), nil, "Should not be nil if git is initialized")
	test.AssertEqual(t, G().Url, "http://localhost:8080", "Should be equal to the remote url")
	test.AssertEqual(t, G().AuthToken, "", "Should be equal to the auth token")
	git = nil
	fatalCalled := false
	fatal = func(a ...interface{}) *pterm.TextPrinter {
		fatalCalled = true
		return nil
	}
	var g *git2.Git
	test.AssertEqual(t, G(), g)
	test.AssertEqual(t, fatalCalled, true, "Should call the fatal function")
}
