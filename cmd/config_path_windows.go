package cmd

import (
	"os"
	"path"
)

var (
	configDir = path.Join(os.Getenv("LOCALAPPDATA"), "containerum")
)
