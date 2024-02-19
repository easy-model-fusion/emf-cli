package command

import "testing"

func TestVersion_runVersion(t *testing.T) {
	runVersion(nil, []string{})
}
