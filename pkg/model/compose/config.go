package compose

type Config struct {
	File     string
	Target   string
	Source   string
	External bool
}

type Configs []Config
