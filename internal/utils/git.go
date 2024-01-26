package utils

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/mod/semver"
)

const GitURL = "https://github.com/SchawnnDev/" // TODO: change when repo is public

// CheckNewSDKVersion checks if a new version of the sdk is available
func CheckNewSDKVersion() bool {
	return false
}

// CheckNewCLIVersion checks if a new version of the cli is available
func CheckNewCLIVersion() bool {
	return false
}

// GetLatestTag returns the latest tag of the given project
func GetLatestTag(project string) (string, error) {
	// Open new remote repository
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{GitURL + project + ".git"},
	})

	// List all references
	refs, err := rem.List(&git.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: git.AppendPeeled,
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
func CloneSDK(tag, to string) error {
	_, err := git.PlainClone(to, false, &git.CloneOptions{
		URL:           GitURL + "sdk.git",
		ReferenceName: plumbing.NewTagReferenceName(tag),
	})

	if err != nil {
		return fmt.Errorf("error cloning sdk: %w", err)
	}
	return nil
}
