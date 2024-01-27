package sdk

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestCheckForUpdates(t *testing.T) {
	dname := test.CreateFullTestSuite(t)
	defer os.RemoveAll(dname)

	err := config.GetViperConfig()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("sdk-tag", "")
	test.AssertEqual(t, checkForUpdates(), false, "Should return false if no tag is set")

	viper.Set("sdk-tag", "v0.0.1")
	test.AssertEqual(t, checkForUpdates(), true, "Should return true if tag is set and there is an update")

	tag, err := utils.GetLatestTag("sdk")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("sdk-tag", tag)
	test.AssertEqual(t, checkForUpdates(), false, "Should return false if tag is set and there is no update")
}

func TestCanSendUpdateSuggestion(t *testing.T) {
	dname := test.CreateFullTestSuite(t)
	defer os.RemoveAll(dname)

	err := config.GetViperConfig()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("update-suggested", false)
	test.AssertEqual(t, canSendUpdateSuggestion(), true, "Should return true if update-suggested is false")

	viper.Set("update-suggested", true)
	test.AssertEqual(t, canSendUpdateSuggestion(), false, "Should return false if update-suggested is true")
}
