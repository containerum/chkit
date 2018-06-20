package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/containerum/chkit/pkg/cli"
	"github.com/containerum/chkit/pkg/model/doc"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gflag"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func main() {
	var flags = struct {
		OutputFile string `desc:"output file in case of defined --commmand flag, -(STDOUT) by default"`
		OutputDir  string `desc:"output directory in case of --hugo"`
		Format     string
		MD         bool   `desc:"generate markdown doc for specific command"`
		Command    string `desc:"command"`
		Hugo       bool
	}{}
	if err := gflag.ParseToDef(&flags); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	flag.Parse()

	var stack = cli.RootCommands()
	var commands = make([]cobra.Command, 0, len(stack))
	for len(stack) > 0 {
		var cmd = pop(&stack)
		if flags.Command != "" && "chkit "+flags.Command == cmd.Name() {
			commands = append(commands, cmd)
			break
		}
		stack = append(stack, cmd.Commands()...)
		commands = append(commands, cmd)
	}
	if flags.Hugo {
		fmt.Println("generating hugo docs")
		var categories = map[string]*bytes.Buffer{}
		for _, cmd := range commands {
			if cmd.Parent() != nil &&
				cmd.Parent().Use != "chkit" {
				var category = cmd.Parent().Use
				if page, ok := categories[category]; !ok || page == nil {
					var categoryDoc = doc.Command{*cmd.Parent()}
					page = bytes.NewBufferString(hugoHeader(category, categoryDoc.Doc()))
					page.WriteString(categoryDoc.Markdown())
					categories[category] = page
				} else {
					page.WriteString(doc.Command{cmd}.Markdown())
				}
				continue
			}
			var category = cmd.Use
			if page, ok := categories[category]; !ok || page == nil {
				var categoryDoc = doc.Command{cmd}
				page = bytes.NewBufferString(hugoHeader(category, categoryDoc.Doc()))
				page.WriteString(categoryDoc.Markdown())
				categories[category] = page
			}
		}
		for category, page := range categories {
			fmt.Println("writing", category)
			ioutil.WriteFile(path.Join(flags.OutputDir, category+".md"), page.Bytes(), os.ModePerm)
		}
		return
	}

	if flags.Command != "" {
		var command = str.Vector(strings.Fields(flags.Command)).Join(" ")
		for _, cmd := range commands {
			if cmd.CommandPath() != command {
				continue
			}
			var docText string
			switch {
			case flags.Format != "":
				var err error
				docText, err = doc.Command{cmd}.Format(flags.Format)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			case flags.MD:
				docText = doc.Command{cmd}.Markdown()
			default:
				docText = doc.Command{cmd}.String()
			}
			if flags.OutputFile == "" {
				fmt.Println(docText)
				return
			}
			if err := ioutil.WriteFile(flags.OutputFile, []byte(docText), os.ModePerm); err != nil {
				fmt.Printf("unable to write doc to file %q:\n\t%v", flags.OutputFile, err)
				os.Exit(1)
			}
		}
		return
	}
}

func pop(stack *[]*cobra.Command) cobra.Command {
	var cmd = (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return *cmd
}

func hugoHeader(category string, cmd doc.Doc) string {
	var header map[string]interface{}
	if category != "" {
		header = map[string]interface{}{
			"title":       strings.Title(cmd.Path),
			"linktitle":   cmd.Path,
			"description": cmd.Description,
			"menu": map[string]interface{}{
				"docs": map[string]interface{}{
					"parent": "commands",
					"weight": 5,
				},
			},
			"weight": 2,
			"draft":  false,
		}
	}
	var data, err = yaml.Marshal(header)
	if err != nil {
		panic(err)
	}
	return "---\n" + string(data) + "\n---\n\n"
}
