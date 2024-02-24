package huggingface

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestAllModulesString(t *testing.T) {
	modules := AllModulesString()

	test.AssertEqual(t, len(modules), len(AllModules), "The number of tags should be equal to the number of AllModules")

	for i, module := range modules {
		test.AssertEqual(t, module, string(AllModules[i]))
	}

}
