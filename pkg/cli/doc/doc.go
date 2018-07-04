package doc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	chkitDoc "github.com/containerum/chkit/pkg/model/doc"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Doc(ctx *context.Context) *cobra.Command {
	var flags struct {
		Output  string `desc:"output file, STDOUT by default"`
		Command string `desc:"print docs for command and its subcommands, example 'chkit doc --command \"create depl\"'"`
		List    bool   `desc:"print command names"`
		Format  string
		MD      bool `desc:"generate markdown docs"`
	}
	var cmd = &cobra.Command{
		Use:   "doc",
		Short: "Print full chkit help",
		Run: func(cmd *cobra.Command, args []string) {
			var currentCommand = cmd.Root()
			var doc = &bytes.Buffer{}
			switch {
			case flags.List:
				for _, command := range getCommandList(currentCommand) {
					if flags.Format == "" {
						doc.WriteString(chkitDoc.Command{*command}.Doc().Path + "\n")
						continue
					}
					var str, err = chkitDoc.Command{*command}.Format(flags.Format)
					if err != nil {
						ferr.Println(err)
						ctx.Exit(1)
					}
					doc.WriteString(str + "\n")
				}
			case flags.Command != "":
				command, _, _ := cmd.Parent().Find(strings.Fields(flags.Command))
				var md = chkitDoc.Command{*command}.Markdown()
				doc.WriteString(md)
			case flags.Command == "":
				for _, command := range getCommandList(currentCommand) {
					command.SetOutput(doc)
					if flags.MD {
						doc.WriteString(chkitDoc.Command{*command}.Markdown())
					} else {
						fmt.Fprintf(doc, "Command : %s\n\n", func() string {
							if command.Parent() != nil && command.Parent().Use != "chkit" {
								return command.Parent().Use + " " + command.Use
							}
							return command.Use
						}())
						//doc.WriteString(getDoc(head))
						command.Usage()
					}
				}
			}
			if flags.Output == "" {
				fmt.Println(doc)
			} else {
				if err := ioutil.WriteFile(flags.Output, doc.Bytes(), os.ModePerm); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}

		},
	}
	if err := gpflag.ParseTo(&flags, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}

func getCommandList(root *cobra.Command) []*cobra.Command {
	var stack = []*cobra.Command{root}
	var commands []*cobra.Command
	for len(stack) > 0 {
		var head = func() *cobra.Command {
			var cmd = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			return cmd
		}()
		commands = append(commands, head)
		stack = append(stack, head.Commands()...)
		commands = append(commands, head.Commands()...)
	}
	return commands
}
