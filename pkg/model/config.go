package model

import (
	"io"
)

type StorableConfig struct {
	UserInfo
	DefaultNamespace string
}
type Config struct {
	StorableConfig
	Tokens
	APIaddr     string
	Fingerprint string
	Log         io.Writer
}
