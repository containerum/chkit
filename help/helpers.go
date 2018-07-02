package help

import (
	"strings"

	"github.com/ninedraft/boxofstuff/str"
	"github.com/spf13/cobra"
)

//go:generate fileb0x b0x.toml
func MustGetString(command string) string {
	data, err := ReadFile(str.Fields(command).
		Map(strings.TrimSpace).
		Join("/") + ".md")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func GetString(command string) (string, error) {
	data, err := ReadFile(str.Fields(command).
		Map(strings.TrimSpace).
		Join("/") + ".md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Command(cmd *cobra.Command) (string, error) {
	var data, err = ReadFile(str.Fields(cmd.CommandPath()).
		Filter(func(str string) bool {
			return str != "chkit"
		}).Join("/") + ".md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Auto(root *cobra.Command) {
	AutoForCommands(root.Commands())
}

func AutoForCommands(cmds []*cobra.Command) {
	var stack = commandStack(cmds)
	var commands []*cobra.Command
	for stack.Len() > 0 {
		var cmd = stack.Pop()
		if cmd == nil {
			panic("NIL CMD")
		}
		stack.Push(cmd.Commands()...)
		commands = append(commands, cmd)
	}
	for _, cmd := range commands {
		if strings.TrimSpace(cmd.Long) == "" {
			var help, err = Command(cmd)
			if err == nil {
				cmd.Long = help
			}
		}
	}
}

type commandStack []*cobra.Command

func (stack *commandStack) Len() int {
	return len(*stack)
}

func (stack *commandStack) Pop() *cobra.Command {
	var head = (*stack)[stack.Len()-1]
	(*stack)[stack.Len()-1] = nil
	*stack = (*stack)[:stack.Len()-1]
	return head
}

func (stack *commandStack) Push(cmds ...*cobra.Command) {
	*stack = append(*stack, cmds...)
}
