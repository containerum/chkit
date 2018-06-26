package help

import (
	"strings"

	"github.com/ninedraft/boxofstuff/str"
)

//go:generate fileb0x b0x.toml
func GetString(command string) string {
	data, err := ReadFile(str.Fields(command).
		Map(strings.TrimSpace).
		Join("/") + ".md")
	if err != nil {
		panic(err)
	}
	return string(data)
}
