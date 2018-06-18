package doc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Doc(ctx *context.Context) *cobra.Command {
	var flags struct {
		Output  string `desc:"output file, STDOUT by default"`
		Command string `desc:"print docs for command and its subcommands, example 'chkit doc --command \"create depl\"'"`
		List    bool   `desc:"print command names"`
		MD      bool   `desc:"generate markdown docs"`
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
					fmt.Fprintf(doc, "%s\n\n", func() string {
						if command.Parent() != nil && command.Parent().Use != "chkit" {
							return command.Parent().Use + " " + command.Use
						}
						return command.Use
					}())
				}
			case flags.Command != "":
				command, _, _ := cmd.Parent().Find(strings.Fields(flags.Command))
				command.SetOutput(doc)
				if flags.MD {
					doc.WriteString(DocMD(command))
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
			case flags.Command == "":
				for _, command := range getCommandList(currentCommand) {
					command.SetOutput(doc)
					if flags.MD {
						doc.WriteString(DocMD(command))
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
					fmt.Println(err)
					os.Exit(1)
				}
			}

		},
	}
	if err := gpflag.ParseTo(&flags, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}

func DocMD(cmd *cobra.Command) string {
	var doc = &bytes.Buffer{}
	fmt.Fprintf(doc, "\n### %s\n\n"+
		"**Aliases**   :\n\n%s\n\n"+
		"**Usage**     :\n\n%s\n\n"+
		"**Example**   :\n\n%s\n\n"+
		"**Flags**     :\n\n%s\n\n"+
		"**Subcommand**:\n\n%s\n\n",
		func() string {
			if cmd.Parent() != nil && cmd.Parent().Use != "chkit" {
				return cmd.Parent().Use + " " + cmd.Use
			}
			return cmd.Use
		}(),
		strings.Join(cmd.Aliases, ", "),
		activekit.OrString(cmd.Long, cmd.Short),
		cmd.Example,
		func() string {
			var d string
			cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
				if !flag.Hidden {
					d += fmt.Sprintf("+ %s %s : %s\n", flag.Name, flag.Shorthand, flag.Usage)
				}
			})
			return text.Indent(d, 2)
		}(),
		func() string {
			var d string
			for _, sub := range cmd.Commands() {
				d += fmt.Sprintf("+ %s : %s\n", sub.Use, sub.Short)
			}
			return text.Indent(d, 2)
		}())
	return doc.String()
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
