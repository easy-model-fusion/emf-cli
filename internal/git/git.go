package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/mod/semver"
	"net/url"
)

type Git struct {
	AuthToken string
	Url       string
}

// NewGit creates a new Git instance
func NewGit(url, authToken string) *Git {
	return &Git{AuthToken: authToken, Url: url}
}

// GenerateAuth generates a new http.BasicAuth
func (g Git) GenerateAuth() *http.BasicAuth {
	if g.AuthToken == "" {
		return nil
	}
	return &http.BasicAuth{Username: "auth", Password: g.AuthToken}
}

// GetProjectUrl returns the url of the given project
func (g Git) GetProjectUrl(project string) (string, error) {
	return url.JoinPath(g.Url, project+".git")
}

// CheckNewSDKVersion checks if a new version of the sdk is available
func (g Git) CheckNewSDKVersion() bool {
	// TODO: implement
	return false
}

// CheckNewCLIVersion checks if a new version of the cli is available
func (g Git) CheckNewCLIVersion() bool {
	// TODO: implement
	return false
}

// GetLatestTag returns the latest tag of the given project
func (g Git) GetLatestTag(project string) (tag string, err error) {
	var remoteUrl string
	if remoteUrl, err = g.GetProjectUrl(project); err != nil {
		return "", fmt.Errorf("get latest tag: %w", err)
	}

	// Open new remote repository
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteUrl},
	})

	var refs []*plumbing.Reference

	// List all references
	refs, err = rem.List(&git.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: git.AppendPeeled,
		Auth:          g.GenerateAuth(),
	})

	if err != nil {
		return "", fmt.Errorf("get latest tag: %w", err)
	}

	// Extract all tags
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	if len(tags) == 0 {
		return "", errors.New("get latest tag: no tags found")
	}

	// Sort tags by semver
	semver.Sort(tags)

	return tags[len(tags)-1], nil
}

// CloneSDK clones the sdk to the given path
func (g Git) CloneSDK(tag, to string) (err error) {
	var remoteUrl string

	if remoteUrl, err = g.GetProjectUrl("sdk"); err != nil {
		return fmt.Errorf("error joining url: %w", err)
	}

	_, err = git.PlainClone(to, false, &git.CloneOptions{
		URL:           remoteUrl,
		ReferenceName: plumbing.NewTagReferenceName(tag),
		Auth:          g.GenerateAuth(),
	})

	if err != nil {
		return fmt.Errorf("error cloning sdk: %w", err)
	}
	return nil
}
