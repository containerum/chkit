package activeingress

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/ninedraft/ranger/intranger"
)

func pathsMenu(services service.ServiceList, paths ingress.PathList) ingress.PathList {
	oldPaths := paths.Copy()
	var ok bool
	for exit := false; !exit; {
		var menu []*activekit.MenuItem
		for ind, path := range paths {
			menu = append(menu, &activekit.MenuItem{
				Label: fmt.Sprintf("Edit path %s", func() string {
					if path.Path != "" && (path.ServicePort >= 0 || path.ServiceName != "") {
						return fmt.Sprintf("%q -> %s:%d", path.Path, path.ServiceName, path.ServicePort)
					} else if path.Path != "" {
						return path.Path
					} else if path.ServicePort >= 0 || path.ServiceName != "" {
						return fmt.Sprintf("%s:%d", path.ServiceName, path.ServicePort)
					}
					return "empty path"
				}()),
				Action: func(ind int) func() error {
					return func() error {
						paths[ind] = editPathMenu(services, paths, ind)
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
						path := ingress.Path{}
						if len(services) == 1 {
							path.ServiceName = services[0].Name
							ports := services[0].AllTargetPorts()
							if len(ports) > 0 {
								path.ServicePort = ports[0]
							}
						}
						ingr, ok := pathMenu(services, path)
						if ok {
							paths = paths.Append(ingr)
						}
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

func editPathMenu(services service.ServiceList, paths ingress.PathList, ind int) ingress.Path {
	path := paths[ind]
	var oldService = path
	var ok bool
	var selectedService *service.Service
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Edit path",
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set path         : %s",
						activekit.OrString(path.Path, "undefined (required)")),
					Action: func() error {
						p := strings.TrimSpace(activekit.Promt("Type path name (hit Enter to leave %s)",
							activekit.OrString(path.Path, "empty")))
						if p == "" {
							return nil
						}
						path.Path = p
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set service      : %s",
						activekit.OrString(path.ServiceName, "undefined (required)")),
					Action: func() error {
						var menu activekit.MenuItems
						for _, serv := range services {
							menu = menu.Append(&activekit.MenuItem{
								Label: serv.Name,
								Action: func(serv service.Service) func() error {
									return func() error {
										path.ServiceName = serv.Name
										cp := serv.Copy()
										selectedService = &cp
										ports := serv.AllExternalPorts()
										if len(ports) > 0 {
											path.ServicePort = ports[0]
										}
										return nil
									}
								}(serv),
							})
						}
						(&activekit.Menu{
							Title: "Select service",
							Items: menu.Append(activekit.MenuItems{
								{
									Label: fmt.Sprintf("Return to previous menu, use %s", path.ServiceName),
								},
							}...),
						}).Run()
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set service port : %s", func() string {
						if path.ServicePort < 0 {
							return "undefined (required)"
						}
						return strconv.Itoa(path.ServicePort)
					}()),
					Action: func() error {
						serv, selectable := func() (service.Service, bool) {
							if selectedService != nil {
								return (*selectedService).Copy(), true
							}
							serv, ok := services.GetByName(path.ServiceName)
							return serv, ok
						}()
						if selectable {
							var menu activekit.MenuItems
							for _, port := range serv.AllExternalPorts() {
								menu = menu.Append(&activekit.MenuItem{
									Label: fmt.Sprintf(" :%d", port),
									Action: func(port int) func() error {
										return func() error {
											path.ServicePort = port
											return nil
										}
									}(port),
								})
							}
							(&activekit.Menu{
								Title: "Select port",
								Items: menu.Append(activekit.MenuItems{
									{
										Label: fmt.Sprintf("Return to previous menu, use port %d", path.ServicePort),
									},
								}...),
							}).Run()
							return nil
						}
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

func pathMenu(services service.ServiceList, path ingress.Path) (ingress.Path, bool) {
	var oldPath = path
	var ok bool
	var selectedService *service.Service
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Edit path",
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set service      : %s",
						activekit.OrString(path.ServiceName, "undefined (required)")),
					Action: func() error {
						var menu activekit.MenuItems
						for _, serv := range services {
							menu = menu.Append(&activekit.MenuItem{
								Label: serv.Name,
								Action: func(serv service.Service) func() error {
									return func() error {
										path.ServiceName = serv.Name
										ports := serv.AllExternalPorts()
										if len(ports) == 1 {
											path.ServicePort = ports[0]
										}
										return nil
									}
								}(serv),
							})
						}
						(&activekit.Menu{
							Title: "Select service",
							Items: menu.Append(activekit.MenuItems{
								{
									Label: fmt.Sprintf("Return to previous menu, use %s", path.ServiceName),
								},
							}...),
						}).Run()
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set service port : %s", func() string {
						if path.ServicePort < 0 {
							return "undefined (required)"
						}
						return strconv.Itoa(path.ServicePort)
					}()),
					Action: func() error {
						serv, selectable := func() (service.Service, bool) {
							if selectedService != nil {
								return (*selectedService).Copy(), true
							}
							serv, ok := services.GetByName(path.ServiceName)
							return serv, ok
						}()
						if selectable {
							var menu activekit.MenuItems
							for _, port := range serv.AllExternalPorts() {
								menu = menu.Append(&activekit.MenuItem{
									Label: fmt.Sprintf(" :%d", port),
									Action: func(port int) func() error {
										return func() error {
											path.ServicePort = port
											return nil
										}
									}(port),
								})
							}
							(&activekit.Menu{
								Title: "Select port",
								Items: menu.Append(activekit.MenuItems{
									{
										Label: fmt.Sprintf("Return to previous menu, use port %d", path.ServicePort),
									},
								}...),
							}).Run()
							return nil
						}
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
					Label: fmt.Sprintf("Set path         : %s",
						activekit.OrString(path.Path, "undefined (required)")),
					Action: func() error {
						p := strings.TrimSpace(activekit.Promt("Type path (hit Enter to leave %s): ",
							activekit.OrString(path.Path, "empty")))
						if p == "" {
							return nil
						}
						path.Path = p
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
		return oldPath, ok
	}
	return path, ok
}
