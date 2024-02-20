package app

import (
	git2 "github.com/easy-model-fusion/emf-cli/internal/git"
	"github.com/easy-model-fusion/emf-cli/internal/huggingface"
	"github.com/pterm/pterm"
)

var huggingFace *huggingface.HuggingFace
var git *git2.Git
var fatal = pterm.Fatal.Println // make it a variable, so we can mock it in tests

func H() *huggingface.HuggingFace {
	if huggingFace == nil {
		fatal("HuggingFace is not initialized, please run InitHuggingFace() first.")
	}
	return huggingFace
}
func G() *git2.Git {
	if git == nil {
		fatal("Git is not initialized, please run InitGit() first.")
	}
	return git

}

// InitHuggingFace Initialize HuggingFace
func InitHuggingFace(baseUrl, proxyUrl string) {
	huggingFace = huggingface.NewHuggingFace(baseUrl, proxyUrl)
}

// InitGit Initialize git
func InitGit(url, authToken string) {
	git = git2.NewGit(url, authToken)
}
