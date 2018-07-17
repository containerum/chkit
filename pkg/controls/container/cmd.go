package container

import (
	"sort"
	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/ninedraft/boxofstuff/str"
)

func componentCmds(cont *container.Container) *activekit.MenuItem {
	var cmds = append([]string{}, cont.Commands...)
	return &activekit.MenuItem{
		Label: "Edit command",
		Action: func() error {
			for exit := false; !exit; {
				cmds = deleteEmpty(cmds)
				sort.Slice(cmds, func(i, j int) bool {
					return cmds[i] < cmds[j]
				})
				var menuCmds activekit.MenuItems
				for i := range cmds {
					menuCmds = append(menuCmds, componentCmd(&cmds[i]))
				}
				menuCmds = menuCmds.Append(&activekit.MenuItem{
					Label: "Add command",
					Action: func() error {
						var cmd string
						componentCmd(&cmd).Action()
						cmds = append(cmds, cmd)
						return nil
					},
				},
					&activekit.MenuItem{
						Label: "Confirm",
						Action: func() error {
							cont.Commands = cmds
							exit = true
							return nil
						},
					},
					&activekit.MenuItem{
						Label: "Drop all changes, return to previous menu",
						Action: func() error {
							exit = true
							return nil
						},
					})
				(&activekit.Menu{
					Title: "Container -> Commands",
					Items: menuCmds,
				}).Run()
			}
			return nil
		},
	}
}

func componentCmd(oldCmd *string) *activekit.MenuItem {
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
