package help

import (
	"strings"

	"github.com/ninedraft/boxofstuff/str"
)

func GetString(command string) string {
	data, err := ReadFile(str.Fields(command).
		Map(strings.TrimSpace).
		Join("/") + ".md")
	if err != nil {
		panic(err)
	}
	return string(data)
}
