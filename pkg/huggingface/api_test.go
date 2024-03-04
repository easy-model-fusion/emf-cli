package huggingface

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestNewHuggingFace(t *testing.T) {
	hug := NewHuggingFace("http://localhost:8080", "")
	test.AssertNotEqual(t, hug, nil, "Should not be nil if huggingface is initialized")
	hugcast, ok := hug.(*huggingFace)
	test.AssertEqual(t, ok, true, "Should be able to cast to huggingFace")
	test.AssertEqual(t, hugcast.BaseUrl, "http://localhost:8080", "Should be equal to the base url")

	// test with proxy
	hug = NewHuggingFace("http://localhost:8080", "http://test:test@localhost:8080")
	test.AssertNotEqual(t, hug, nil, "Should not be nil if huggingface is initialized")
	hugcast, ok = hug.(*huggingFace)
	test.AssertEqual(t, ok, true, "Should be able to cast to huggingFace")
	test.AssertEqual(t, hugcast.BaseUrl, "http://localhost:8080", "Should be equal to the base url")
	test.AssertNotEqual(t, hugcast.Client, nil, "Should not be nil")
}
