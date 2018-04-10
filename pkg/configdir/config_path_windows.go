package configdir

import (
	"path"
)

var (
	configDir = path.Join("AppData", "Local", "containerum")
	logDir    = path.Join(configDir, "support")
)
