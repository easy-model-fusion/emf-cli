package utils

import (
	"github.com/easy-model-fusion/client/test"
	"runtime"
	"testing"
)

func TestCheckForExecutable(t *testing.T) {
	test.AssertEqual(t, CheckForExecutable("anexecutablethatcouldnotexists-yeahhh"), true)

	switch runtime.GOOS {
	case "windows":
		test.AssertEqual(t, CheckForExecutable("cmd.exe"), true)
	case "darwin":
		fallthrough
	case "linux":
		test.AssertEqual(t, CheckForExecutable("ls"), true)
	default:
		t.Skip("Unsupported operating system")
	}
}
