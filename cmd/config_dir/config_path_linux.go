package confDir

import "path"

var (
	configDir = path.Join(".config", "containerum")
	logDir    = path.Join(configDir, "support")
)
