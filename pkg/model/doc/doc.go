package doc

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	ListFormat       = "{{.Path}}"
	MarkdownTemplate = `#### <a name="{{.Link}}">{{.Path}}</a>` + "\n\n" +
		"**Description**:\n\n{{.Description}}\n\n" +
		"**Example**:\n\n{{.Example}}\n\n" +
		"**Flags**:\n\n" +
		"| Short | Name | Usage | Default value |\n" +
		"| ----- | ---- | ----- | ------------- |\n" +
		"{{range .Flags}}" +
		"| {{if .Short}}-{{.Short}}{{end}} " +
		"| --{{.Name}} " +
		"| {{.Usage}} " +
		"| {{if ne .DefValue \"[]\"}}{{.DefValue}}{{end}} " +
		"|\n" +
		"{{end}}\n\n" +
		"**Subcommands**:\n\n" +
		"{{range .Subcommands}}" +
		"* **[{{.Name}}](#{{.Link}})** {{.ShortDescription}}\n" +
		"{{end}}\n\n"
	TextTamplate = "Command: {{.Path}}\n" +
		"Description:\n{{.Description}}\n" +
		"Example:\n{{.Example}}\n" +
		"Flags:\n" +
		"{{range .Flags}}" +
		"{{if .Short}}-{{.Short}}{{else}}  {{end}} " +
		"--{{.Name}} " +
		"{{.Usage}} " +
		"{{if ne .DefValue \"[]\"}}{{.DefValue}}{{end}} " +
		"\n" +
		"{{end}}\n"
)

type Command struct {
	cobra.Command
}

type Doc struct {
	Link        string        `json:"link"`
	Path        string        `json:"path"`
	Name        string        `json:"name"`
	Description string        `json:"help"`
	Example     string        `json:"example"`
	Flags       []sflags.Flag `json:"flags"`
	Subcommands []SubCommand  `json:"subcommands,omitempty"`
}

type SubCommand struct {
	Link             string `json:"link"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
}

func (cmd Command) String() string {
	var str, err = cmd.Format(TextTamplate)
	if err != nil {
		panic(fmt.Errorf("Command.String: %v", err))
	}
	return str
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
	var subc = make([]SubCommand, 0, len(cmd.Commands()))
	for _, sc := range cmd.Commands() {
		subc = append(subc, SubCommand{
			Link: str.Vector{sc.CommandPath()}.
				Map(str.TrimPrefix("chkit ")).
				Map(strings.NewReplacer(" ", "_").Replace).
				FirstNonEmpty(),
			Name:             strings.TrimPrefix(sc.CommandPath(), "chkit "),
			ShortDescription: sc.Short,
		})
	}
	var cmdPath = strings.TrimPrefix(cmd.CommandPath(), "chkit ")
	return Doc{
		Link: strings.Replace(cmdPath, " ", "_", -1),
		Path: cmdPath,
		Name: cmd.Name(),
		Description: str.Vector{
			cmd.Long,
			cmd.Short,
			cmd.Example,
		}.FirstNonEmpty(),
		Example:     cmd.Example,
		Flags:       cmdFlags,
		Subcommands: subc,
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
		panic(fmt.Errorf("Command.Markdown: %v", err))
	}
	return str
}
