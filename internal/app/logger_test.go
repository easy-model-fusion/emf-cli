package app

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestInitLogger(t *testing.T) {
	initLogger()
	test.AssertNotEqual(t, logger, nil, "Should not be nil if logger is initialized")
	test.AssertNotEqual(t, L(), nil, "Should not be nil if logger is initialized")
}
