package main

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
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
	models, err := app.GetDownloadedModelNames()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(models)
}
