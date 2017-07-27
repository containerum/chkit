package dbconfig

import "time"

var (
	DefaultTCPServer  string
	DefaultHTTPServer string
)

const (
	DefaultBufferSize  = 1024
	DefaultHTTPTimeout = 10 * time.Second
	DefaultTCPTimeout  = 10 * time.Second
)
