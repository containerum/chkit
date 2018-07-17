package deplactive

import (
	"strings"

	"sort"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/ninedraft/boxofstuff/str"
)

func componentCmd(cont *container.Container) activekit.MenuItems {
	cont.Commands = deleteEmpty(cont.Commands)
	sort.Slice(cont.Commands, func(i, j int) bool {
		return cont.Commands[i] < cont.Commands[j]
	})
	var items activekit.MenuItems
	for i := range cont.Commands {
		items = items.Append(componentEditCmd(&cont.Commands[i]))
	}
	return items.Append(&activekit.MenuItem{
		Label: "Add command",
		Action: func() error {
			var cmd string
			if componentEditCmd(&cmd).Action() == nil {
				cont.Commands = append(cont.Commands, cmd)
			}
			return nil
		},
	})
}

func componentEditCmd(oldCmd *string) *activekit.MenuItem {
	var label = str.Vector{*oldCmd, "empty cmd"}.FirstNonEmpty()
	var cmd = *oldCmd
	return &activekit.MenuItem{
		Label: "Edit command " + label,
		Action: func() error {
			cmd = activekit.Promt("Type command, hit Enter to save or leave empty to delete: ")
			cmd = strings.TrimSpace(cmd)
			*oldCmd = cmd
			return nil
		},
	}
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
