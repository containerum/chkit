package deplactive

import (
	"fmt"
	"strings"

	"strconv"

	"io"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/containerum/kube-client/pkg/model"
)

func componentEditContainers(config Wizard) activekit.MenuItems {
	var menuItems = make(activekit.MenuItems, 0, len(config.Deployment.Containers))
	for _, container := range config.Deployment.Containers {
		menuItems = menuItems.Append(componentEditContainer(
			&container,
			config.Volumes,
			config.Configmaps))
	}
	return menuItems.Append(&activekit.MenuItem{
		Label: "Create container",
		Action: func() error {
			var cont = container.Container{
				Container: model.Container{
					Limits: model.Resource{
						CPU:    20 * MIN_CPU,
						Memory: 20 * MIN_MEM,
					},
					Name: namegen.Aster(),
				},
			}
			if componentEditContainer(&cont,
				config.Volumes, config.Configmaps).Action() == nil {
				config.Deployment.Containers = append(config.Deployment.Containers, cont)
			}
			return nil
		},
	})
}

func componentEditContainer(oldCont *container.Container, volumes, configmaps []string) *activekit.MenuItem {
	var cont = func() *container.Container {
		var c = oldCont.Copy()
		return &c
	}()
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit container %s [%s]", cont.Name,
			activekit.OrString(cont.Image, "undefined image")),
		Action: func() error {
			for exit := false; !exit; {
				_, err := (&activekit.Menu{
					Title: "Deployment -> Container",
					Items: activekit.MenuItems{
						componentEditContainerName(cont),
						componentEditContainerImage(cont),
					}.
						Append(componentEnvs(cont)...).
						Append(
							componentEditCPU(cont),
							componentEditMemory(cont),
						).
						Append(componentEditContainerVolumes(cont, volumes)...).
						Append(componentEditContainerConfigmaps(cont, configmaps)...).
						Append(
							&activekit.MenuItem{
								Label: "Confirm",
								Action: func() error {
									if err := ValidateContainer(*cont); err != nil {
										fmt.Println(err)
									} else {
										exit = true
										*oldCont = *cont
									}
									return nil
								},
							},
							&activekit.MenuItem{
								Label: "Drop changes, return to previous menu",
								Action: func() error {
									exit = true
									return io.EOF
								},
							}),
				}).Run()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func componentEditContainerName(cont *container.Container) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit name : %s",
			activekit.OrString(cont.Name, "undefined, required")),
		Action: func() error {
			for {
				var name = activekit.Promt("Type container name (hit Enter to leave %s): ",
					activekit.OrString(cont.Name, "empty"))
				name = strings.TrimSpace(name)
				if err := validation.ValidateLabel(name); name != "" && err == nil {
					cont.Name = name
				} else if name != "" && err != nil {
					fmt.Printf("%s is invalid container name\n", name)
					continue
				}
				break
			}
			return nil
		},
	}
}

func componentEditContainerImage(cont *container.Container) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set image : %s",
			activekit.OrString(cont.Image, "undefined, required")),
		Action: func() error {
			var image = activekit.Promt("Type image (hit Enter to leave %s): ",
				activekit.OrString(cont.Image, "empty"))
			if err := validation.ValidateImageName(image); image != "" && err == nil {
				cont.Image = image
			} else if err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}
}

func componentEditCPU(cont *container.Container) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set CPU limit: %d mCPU", cont.Limits.CPU),
		Action: func() error {
			for {
				var cpuStr = activekit.Promt("Type mCPU (hit Enter to use %d mCPU, expected value in %v mCPU): ",
					cont.Limits.CPU, CPULimit)
				cpuStr = strings.TrimSpace(cpuStr)
				if cpu, err := strconv.ParseUint(cpuStr, 10, 32); cpuStr != "" && err == nil {
					if !CPULimit.Containing(int(cpu)) {
						fmt.Printf("CPU limit must be number in %v\n", CPULimit)
						continue
					}
					cont.Limits.CPU = uint(cpu)
				} else if err != nil {
					fmt.Printf("%q is invalid CPU limit\n", cpuStr)
					continue
				}
				break
			}
			return nil
		},
	}
}

func componentEditMemory(cont *container.Container) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set memory limit: %d Mb", cont.Limits.Memory),
		Action: func() error {
			for {
				var memoryStr = activekit.Promt("Type memory limit (hit Enter to use %d Mb, expected value in %v Mb): ",
					cont.Limits.Memory, MemLimit)
				memoryStr = strings.TrimSpace(memoryStr)
				if memory, err := strconv.ParseUint(memoryStr, 10, 32); memoryStr != "" && err == nil {
					if !MemLimit.Containing(int(memory)) {
						fmt.Printf("Memory limit number must be number in %v\n", MemLimit)
						continue
					}
					cont.Limits.Memory = uint(memory)
				} else if err != nil {
					fmt.Printf("%q is invalid memory limit\n", memoryStr)
					continue
				}
				break
			}
			return nil
		},
	}
}
