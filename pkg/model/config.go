package model

import (
	"io"
)

type StorableConfig struct {
	UserInfo
}
type Config struct {
	StorableConfig
	Tokens
	APIaddr     string
	Fingerprint string
	Log         io.Writer
}
