package confDir

import (
	"os"
	"path"
)

var (
	ConfigDir = path.Join(os.Getenv("LOCALAPPDATA"), "containerum")
)
