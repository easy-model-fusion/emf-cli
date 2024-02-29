package app

import (
	"github.com/easy-model-fusion/emf-cli/internal/git"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
)

var huggingFace huggingface.HuggingFace
var gitInstance *git.Git
var fatal = pterm.Fatal.Println // make it a variable, so we can mock it in tests

func H() huggingface.HuggingFace {
	if huggingFace == nil {
		fatal("HuggingFace is not initialized, please run InitHuggingFace() first.")
	}
	return huggingFace
}
func G() *git.Git {
	if gitInstance == nil {
		fatal("Git is not initialized, please run InitGit() first.")
	}
	return gitInstance

}

// InitHuggingFace Initialize HuggingFace
func InitHuggingFace(baseUrl, proxyUrl string) {
	huggingFace = huggingface.NewHuggingFace(baseUrl, proxyUrl)
}

// InitGit Initialize git
func InitGit(url, authToken string) {
	gitInstance = git.NewGit(url, authToken)
}
