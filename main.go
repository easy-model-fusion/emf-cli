package main

import (
	"github.com/easy-model-fusion/emf-cli/cmd"
	"github.com/easy-model-fusion/emf-cli/internal/app"
)

var (
	// Version is the binary version + build number
	Version = ""
	// BuildDate is the date of build
	BuildDate = ""
)

func main() {
	app.Init(Version, BuildDate)

	// Execute command
	cmd.Execute()
}
