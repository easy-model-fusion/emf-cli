package utils

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestCheckAskForPython_Success(t *testing.T) {
	// check python
	a, ok := CheckForPython()
	if !ok {
		return
	}

	b, ok := CheckAskForPython()
	test.AssertEqual(t, ok, true, "Should return true if python is installed")
	test.AssertEqual(t, a, b, "Should return the same value as CheckForPython")
}
