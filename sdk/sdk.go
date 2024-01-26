package sdk

import "embed"

//go:embed main.py config.yaml .gitignore
var EmbeddedFiles embed.FS
