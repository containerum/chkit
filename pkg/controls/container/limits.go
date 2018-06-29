package container

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/limits"
	"github.com/containerum/chkit/pkg/util/activekit"
)

func componentLimits(cont *container.Container) *activekit.MenuItem {
	var item = &activekit.MenuItem{
		Label: "Edit limits",
		Action: func() error {
			var contLimits = cont.Limits
			for exit := false; !exit; {
				(&activekit.Menu{
					Title: "Container -> Limits",
					Items: activekit.MenuItems{
						{
							Label: "Set CPU limit",
							Action: func() error {
								for {
									var cpuStr = activekit.Promt("Type CPU limit, hit Enter to leave %d, expected %v mCPU: ", contLimits.CPU, limits.CPULimit)
									cpuStr = strings.TrimSpace(cpuStr)
									if cpu, err := strconv.ParseUint(cpuStr, 10, 32); cpuStr != "" && err != nil {
										fmt.Printf("invalid input %q: %v", cpuStr, err)
										continue
									} else if !limits.CPULimit.Containing(int(cpu)) {
										fmt.Printf("expected limit %v mCPU, got %d", limits.CPULimit, cpu)
										continue
									} else {
										contLimits.CPU = uint(cpu)
									}
									return nil
								}
							},
						},
						{
							Label: "Set memory limit",
							Action: func() error {
								for {
									var memStr = activekit.Promt("Type memory limit, hit Enter to leave %d, expected %v Mb: ", contLimits.Memory, limits.MemLimit)
									memStr = strings.TrimSpace(memStr)
									if memory, err := strconv.ParseUint(memStr, 10, 32); memStr != "" && err != nil {
										fmt.Printf("invalid input %q: %v", memStr, err)
										continue
									} else if !limits.CPULimit.Containing(int(memory)) {
										fmt.Printf("expected limit %v Mb, got %d", limits.MemLimit, memory)
										continue
									} else {
										contLimits.Memory = uint(memory)
									}
									return nil
								}
							},
						},
						{
							Label: "Confirm",
							Action: func() error {
								cont.Limits = contLimits
								exit = true
								return nil
							},
						},
						{
							Label: "Drop all changes and return to previous menu",
							Action: func() error {
								exit = true
								return nil
							},
						},
					},
				}).Run()
			}
			return nil
		},
	}
	return item
}
