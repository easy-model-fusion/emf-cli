package app

import "github.com/pterm/pterm"

var logger *pterm.Logger

func L() *pterm.Logger {
	return logger
}

func initLogger() {
	logger = pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace).WithTime(false)
}
