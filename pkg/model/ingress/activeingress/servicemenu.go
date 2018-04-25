package activeingress

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/ninedraft/ranger/intranger"
)

func pathsMenu(paths ingress.PathList) ingress.PathList {
	oldPaths := paths.Copy()
	var ok bool
	for exit := false; !exit; {
		var menu []*activekit.MenuItem
		for ind, path := range paths {
			menu = append(menu, &activekit.MenuItem{
				Label: fmt.Sprintf("Edit %q", path.ServiceName),
				Action: func(ind int) func() error {
					return func() error {
						paths[ind] = editPathMenu(paths, ind)
						return nil
					}
				}(ind),
			})
		}
		(&activekit.Menu{
			Title: "Edit paths",
			Items: append(menu, []*activekit.MenuItem{
				{
					Label: "Add path",
					Action: func() error {
						paths = paths.Append(pathMenu(ingress.Path{}))
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						exit = true
						ok = true
						return nil
					},
				},
				{
					Label: "Return to previous menu, discard all changes",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			}...),
		}).Run()
	}
	if !ok {
		return oldPaths
	}
	return paths
}

func editPathMenu(paths ingress.PathList, ind int) ingress.Path {
	path := paths[ind]
	var oldService = path
	var ok bool
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Edit path",
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set path      :  %s",
						activekit.OrString(path.ServiceName, "undefined (required)")),
					Action: func() error {
						servName := strings.TrimSpace(activekit.Promt("Type path name (hit Enter to leave %s)",
							activekit.OrString(path.ServiceName, "empty")))
						if servName == "" {
							return nil
						}
						path.ServiceName = servName
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set path port : %s", func() string {
						if path.ServicePort < 0 {
							return "undefined (required)"
						}
						return strconv.Itoa(path.ServicePort)
					}()),
					Action: func() error {
						portLimits := intranger.IntRanger(1, 65553)
						portString := strings.TrimSpace(activekit.Promt("Type port (%v, hit Enter to leave %d)", portLimits, path.ServicePort))
						if portString == "" {
							return nil
						}
						if port, err := strconv.Atoi(portString); err != nil || !portLimits.Containing(port) {
							fmt.Printf("Expect number %v, got %d\n", portLimits, port)
						} else {
							path.ServicePort = port
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						if err := ValidatePath(path); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						exit = true
						ok = true
						return nil
					},
				},
				{
					Label: "Delete path",
					Action: func() error {
						if activekit.YesNo("Are you sure you want to delete path?") {
							paths.Delete(ind)
							exit = true
							ok = true
						}
						return nil
					},
				},
				{
					Label: "Return to previous menu, discard all changes",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			},
		}).Run()
	}
	if !ok {
		return oldService
	}
	return path
}

func pathMenu(path ingress.Path) ingress.Path {
	var oldService = path
	var ok bool
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Edit path",
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set path      :  %s",
						activekit.OrString(path.ServiceName, "undefined (required)")),
					Action: func() error {
						servName := strings.TrimSpace(activekit.Promt("Type path name (hit Enter to leave %s)",
							activekit.OrString(path.ServiceName, "empty")))
						if servName == "" {
							return nil
						}
						path.ServiceName = servName
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set path port : %s", func() string {
						if path.ServicePort < 0 {
							return "undefined (required)"
						}
						return strconv.Itoa(path.ServicePort)
					}()),
					Action: func() error {
						portLimits := intranger.IntRanger(1, 65553)
						portString := strings.TrimSpace(activekit.Promt("Type port (%v, hit Enter to leave %d)", portLimits, path.ServicePort))
						if portString == "" {
							return nil
						}
						if port, err := strconv.Atoi(portString); err != nil || !portLimits.Containing(port) {
							fmt.Printf("Expect number %v, got %d\n", portLimits, port)
						} else {
							path.ServicePort = port
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						if err := ValidatePath(path); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						exit = true
						ok = true
						return nil
					},
				},
				{
					Label: "Return to previous menu, discard all changes",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			},
		}).Run()
	}
	if !ok {
		return oldService
	}
	return path
}
