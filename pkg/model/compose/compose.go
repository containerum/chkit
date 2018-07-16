package compose

type Compose struct {
	Version  string
	Name     string `yaml:"-"`
	Services Services
	Secrets  Configs
	Configs  Configs
}
