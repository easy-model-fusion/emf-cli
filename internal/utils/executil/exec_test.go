package executil

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestCheckForExecutable(t *testing.T) {
	_, ok := CheckForExecutable("anexecutablethatcouldnotexists-yeahhh")
	test.AssertEqual(t, ok, false)
	//switch runtime.GOOS {
	//case "windows":
	//	_, ok = CheckForExecutable("cmd.exe")
	//	test.AssertNotEqual(t, ok, true)
	//case "darwin":
	//	fallthrough
	//case "linux":
	//	_, ok = CheckForExecutable("ls")
	//	test.AssertNotEqual(t, ok, true)
	//default:
	//	t.Skip("Unsupported operating system")
	//}
}
