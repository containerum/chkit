package model

type StorableConfig struct {
	UserInfo
}
type Config struct {
	StorableConfig
	Tokens
	APIaddr     string
	Fingerprint string
}
