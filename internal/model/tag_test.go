package model

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestAllTagsString(t *testing.T) {
	tags := AllTagsString()

	test.AssertEqual(t, len(tags), len(AllTags), "The number of tags should be equal to the number of AllTags")

	for i, tag := range tags {
		test.AssertEqual(t, tag, string(AllTags[i]), "The tag should be equal to TextToText")
	}

}
