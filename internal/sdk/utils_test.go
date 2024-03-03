package sdk

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.InitGit("https://github.com/SchawnnDev", "")
	os.Exit(m.Run())
}

func TestCheckForUpdates(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("sdk-tag", "")
	_, ok := checkForUpdates()
	test.AssertEqual(t, ok, false, "Should return false if no tag is set")

	viper.Set("sdk-tag", "v0.0.1")
	_, ok = checkForUpdates()
	test.AssertEqual(t, ok, true, "Should return true if tag is set and there is an update")

	tag, err := app.G().GetLatestTag("sdk")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("sdk-tag", tag)
	_, ok = checkForUpdates()
	test.AssertEqual(t, ok, false, "Should return false if tag is set and there is no update")
}

func TestCanSendUpdateSuggestion(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, canSendUpdateSuggestion(), true, "Should return true if update-suggested is not set")

	viper.Set("update-suggested", false)
	test.AssertEqual(t, canSendUpdateSuggestion(), true, "Should return true if update-suggested is false")

	viper.Set("update-suggested", true)
	test.AssertEqual(t, canSendUpdateSuggestion(), false, "Should return false if update-suggested is true")
}

func TestResetUpdateSuggestion(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("update-suggested", true)
	ResetUpdateSuggestion()
	test.AssertEqual(t, viper.GetBool("update-suggested"), false, "Should set update-suggested to false")
}

func TestSetUpdateSuggestion(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	setUpdateSuggestion(true)
	test.AssertEqual(t, viper.GetBool("update-suggested"), true, "Should set update-suggested to true")

	setUpdateSuggestion(false)
	test.AssertEqual(t, viper.GetBool("update-suggested"), false, "Should set update-suggested to false")
}

func TestSendUpdateSuggestion(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("update-suggested", false)
	viper.Set("sdk-tag", "")
	SendUpdateSuggestion()
	test.AssertEqual(t, viper.GetBool("update-suggested"), false, "Should not set update-suggested to true if there is no tag")

	viper.Set("update-suggested", false)
	viper.Set("sdk-tag", "v0.0.1")
	SendUpdateSuggestion()
	test.AssertEqual(t, viper.GetBool("update-suggested"), true, "Should set update-suggested to true if there is a tag and update-suggested is false")

	viper.Set("update-suggested", true)
	viper.Set("sdk-tag", "v0.0.1")
	SendUpdateSuggestion()
	test.AssertEqual(t, viper.GetBool("update-suggested"), true, "Should not set update-suggested to true if there is a tag and update-suggested is true")
}

func TestUpgrade(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	app.SetUI(&test.MockUI{})

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("sdk-tag", "")
	err = Upgrade()
	test.AssertNotEqual(t, err, nil, "Should return an error if no tag is set")

	tag, err := app.G().GetLatestTag("sdk")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viper.Set("sdk-tag", tag)
	err = Upgrade()
	test.AssertNotEqual(t, err, nil, "Should return an error if tag is set and there is no update")

	viper.Set("sdk-tag", "v0.0.1")
	err = Upgrade()
	if err != nil {
		t.Error(err)
	}
	test.AssertEqual(t, err, nil, "Should not return an error if tag is set and there is an update")

	// check if config was written with the new tag
	test.AssertEqual(t, viper.GetString("sdk-tag"), tag, "Should write the new tag to the config")

}
