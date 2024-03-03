package test

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"path/filepath"
)

type MockGit struct {
	Tag            string
	LatestTagError error
	CloneSDKError  error
}

func (g *MockGit) GenerateAuth() *http.BasicAuth {
	return nil
}

func (g *MockGit) GetAuthToken() *string {
	return nil
}

func (g *MockGit) GetUrl() *string {
	return nil
}

func (g *MockGit) GetProjectUrl(project string) (string, error) {
	return "", nil
}

func (g *MockGit) CheckNewSDKVersion() bool {
	return false
}

func (g *MockGit) CheckNewCLIVersion() bool {
	return false
}

func (g *MockGit) GetLatestTag(project string) (tag string, err error) {
	return g.Tag, g.LatestTagError
}

func (g *MockGit) CloneSDK(tag, to string) (err error) {
	// create ".git" folder in to
	_ = os.MkdirAll(filepath.Join(to, ".git"), os.ModePerm)
	return g.CloneSDKError
}
