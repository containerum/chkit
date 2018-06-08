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
)

func Doc(ctx *context.Context) *cobra.Command {
	var flags struct {
		Output  string `desc:"output file, STDOUT by default"`
		Command string `desc:"print docs for command and its subcommands, example 'chkit doc --command \"create depl\"'"`
		MD      bool   `desc:"generate markdown docs"`
	}
	var cmd = &cobra.Command{
		Use:   "doc",
		Short: "Print full chkit help",
		Run: func(cmd *cobra.Command, args []string) {
			var currentCommand *cobra.Command
			if flags.Command == "" {
				currentCommand = cmd.Root()
			} else {
				currentCommand, _, _ = cmd.Parent().Find(strings.Fields(flags.Command))
			}
			var doc = &bytes.Buffer{}
			var stack = []*cobra.Command{currentCommand}
			for len(stack) > 0 {
				var head = func() *cobra.Command {
					var cmd = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					return cmd
				}()
				head.SetOutput(doc)
				stack = append(stack, head.Commands()...)
				if flags.MD {
					doc.WriteString(DocMD(head))
				} else {
					fmt.Fprintf(doc, "\n------------------------\nCommand : %s\n", func() string {
						if head.Parent() != nil && head.Parent().Use != "chkit" {
							return head.Parent().Use + " " + head.Use
						}
						return head.Use
					}())

					//doc.WriteString(getDoc(head))
					head.Usage()
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
	fmt.Fprintf(doc, "\n##%s\n\n"+
		"###Aliases:\n  %s\n"+
		"###Usage  :\n %s\n"+
		"###Example:\n  %s\n"+
		"###Flags  :\n%s\n"+
		"###Subcommands :\n%s\n",
		func() string {
			if cmd.Parent() != nil && cmd.Parent().Use != "chkit" {
				return cmd.Parent().Use + " " + cmd.Use
			}
			return cmd.Use
		}(),
		strings.Join(cmd.Aliases, ", "),
		activekit.OrString(cmd.Long, cmd.Short),
		cmd.Example,
		cmd.LocalFlags().FlagUsages(),
		func() string {
			var d string
			for _, sub := range cmd.Commands() {
				d += fmt.Sprintf("+ %s : %s\n", sub.Use, sub.Short)
			}
			return text.Indent(d, 2)
		}())
	return doc.String()
}
