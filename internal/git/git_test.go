package git

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"testing"
)

var g Git

const gitTestUrl = "https://github.com/SchawnnDev"

func TestMain(m *testing.M) {
	g = NewGit(gitTestUrl, "")
	os.Exit(m.Run())
}

func TestNewGit(t *testing.T) {
	// Test with a valid url
	test.AssertNotEqual(t, g, nil, "Expected the git object to be not nil")
	test.AssertEqual(t, *g.GetUrl(), gitTestUrl)
	test.AssertEqual(t, *g.GetAuthToken(), "")
}

func TestGit_GenerateAuth(t *testing.T) {
	var httpNil *http.BasicAuth

	// Test with no auth token
	test.AssertEqual(t, g.GenerateAuth(), httpNil, "Expected the auth token to be nil")

	// Test with an invalid auth token
	n := NewGit(gitTestUrl, "SchawnnDev")
	auth := n.GenerateAuth()
	test.AssertNotEqual(t, auth, httpNil, "Expected the auth token to be not nil")
	test.AssertEqual(t, auth.Password, "SchawnnDev")
}

func TestGit_GetProjectUrl(t *testing.T) {
	// Test with a valid project
	url, err := g.GetProjectUrl("emf-cli")
	test.AssertEqual(t, err, nil, "Expected no error")
	test.AssertEqual(t, url, gitTestUrl+"/emf-cli.git")

	// Test with an invalid project
	n := NewGit("%%%QS%D smd:m invalid", "")
	_, err = n.GetProjectUrl("invalid")
	test.AssertNotEqual(t, err, nil, "Expected error")
}

func TestGit_CheckNewSDKVersion(t *testing.T) {
	// Test with no implementation
	test.AssertEqual(t, g.CheckNewSDKVersion(), false, "Expected false")
}

func TestGit_CheckNewCLIVersion(t *testing.T) {
	// Test with no implementation
	test.AssertEqual(t, g.CheckNewCLIVersion(), false, "Expected false")
}

func TestCloneSDK(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// no sdk was found
	*g.GetUrl() = "https://github.com/SchawnnDevArchive"
	err = g.CloneSDK("v1.0.0", tmpDir)
	test.AssertNotEqual(t, err, nil, "Expected error")

	// check invalid
	*g.GetUrl() = "invalid %%%)à invalid"
	err = g.CloneSDK("v1.0.0", tmpDir)
	test.AssertNotEqual(t, err, nil, "Expected error")
	*g.GetUrl() = gitTestUrl

	// Clone the SDK
	err = g.CloneSDK("v1.0.0", tmpDir)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Check if the SDK was cloned (.git folder exists)
	_, err = os.Stat(fileutil.PathJoin(tmpDir, ".git"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestGetLatestTag(t *testing.T) {
	// check invalid
	*g.GetUrl() = "invalid %%%)à invalid"
	_, err := g.GetLatestTag("sdk")
	test.AssertNotEqual(t, err, nil, "Expected error")
	*g.GetUrl() = gitTestUrl

	// Test with a valid tag
	tag, err := g.GetLatestTag("sdk")
	test.AssertEqual(t, err, nil, "Expected no error")
	t.Log(tag)

	// Test with an invalid project
	_, err = g.GetLatestTag("invalid")
	test.AssertNotEqual(t, err, nil, "Expected error")

	// Test with no tags
	_, err = g.GetLatestTag("SoS")
	test.AssertNotEqual(t, err, nil, "Expected error")
}
