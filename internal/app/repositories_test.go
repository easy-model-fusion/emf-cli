package app

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

func TestInitHuggingFace(t *testing.T) {
	InitHuggingFace("http://localhost:8080", "")
	test.AssertNotEqual(t, huggingFace, nil, "Should not be nil if huggingface is initialized")
	test.AssertNotEqual(t, H(), nil, "Should not be nil if huggingface is initialized")
	test.AssertEqual(t, H().BaseUrl, "http://localhost:8080", "Should be equal to the base url")
}
