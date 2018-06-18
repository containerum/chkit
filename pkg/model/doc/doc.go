package doc

import (
	"bytes"
	"text/template"

	"strings"

	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	ListFormat       = "{{.Name}}"
	MarkdownTemplate = "### {{.Path}}\n\n" +
		"**Description**:\n\n{{.Description}}\n\n" +
		"**Example**:\n\n{{.Example}}\n\n" +
		"**Flags**:\n\n" +
		"| Short | Name | Usage | Default value |\n" +
		"| ----- | ---- | ----- | ------------- |\n" +
		"{{range .Flags}}" +
		"| {{if .Short}}-{{.Short}}{{end}} " +
		"| {{.Name}} " +
		"| {{.Usage}} " +
		"| {{if ne .DefValue \"[]\"}}{{.DefValue}}{{end}} " +
		"|\n" +
		"{{end}}\n\n"
)

type Command struct {
	cobra.Command
}

type Doc struct {
	Path        string        `json:"path"`
	Name        string        `json:"name"`
	Description string        `json:"help"`
	Example     string        `json:"example"`
	Flags       []sflags.Flag `json:"flags"`
}

func (cmd Command) Doc() Doc {
	var cmdFlags []sflags.Flag
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			cmdFlags = append(cmdFlags, sflags.Flag{
				Name:     flag.Name,
				Short:    flag.Shorthand,
				Usage:    strings.Replace(flag.Usage, "\n", " ", -1),
				DefValue: flag.DefValue,
			})
		}
	})
	return Doc{
		Path: strings.TrimPrefix(cmd.CommandPath(), "chkit "),
		Name: cmd.Name(),
		Description: strings.Replace(str.Vector{
			cmd.Long,
			cmd.Short,
			cmd.Example,
		}.FirstNonEmpty(), "\n", " ", -1),
		Example: cmd.Example,
		Flags:   cmdFlags,
	}
}

func (cmd Command) Format(format string) (string, error) {
	var templ, err = template.New(cmd.Name()).Parse(format)
	if err != nil {
		return "", err
	}
	var buf = &bytes.Buffer{}
	err = templ.Execute(buf, cmd.Doc())
	return buf.String(), err
}

func (cmd Command) Markdown() string {
	var str, err = cmd.Format(MarkdownTemplate)
	if err != nil {
		panic(err)
	}
	return str
}
