package model

type Config struct {
	ConfigPath string
	Client     ClientConfig
}

type ClientConfig struct {
	APIaddr  string
	Username string
	Password string
}
