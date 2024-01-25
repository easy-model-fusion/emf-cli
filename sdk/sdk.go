package sdk

import "embed"

//go:embed main.py config.yaml
var EmbeddedFiles embed.FS
