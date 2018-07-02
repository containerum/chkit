package help

import (
	"strings"

	"github.com/ninedraft/boxofstuff/str"
	"github.com/spf13/cobra"
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

func Command(cmd *cobra.Command) string {
	var data, err = ReadFile(str.Fields(cmd.CommandPath()).
		Filter(func(str string) bool {
			return str != "chkit"
		}).Join("/") + ".md")
	if err != nil {
		panic(err)
	}
	return string(data)
}
