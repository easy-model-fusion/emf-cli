package sdk

import "embed"

//go:embed main.py config.yaml .gitignore README.md requirements.txt
var EmbeddedFiles embed.FS
