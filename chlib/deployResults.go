package chlib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type deployResult struct {
	Data struct {
		Items []struct {
			Status struct {
				AvailableReplicas int `json:"availableReplicas"`
			} `json:"status"`
			Metadata struct {
				CreatedAt time.Time `json:"creationTimestamp"`
				Name      string    `json:"name"`
			} `json:"metadata"`
			Spec struct {
				Replicas int `json:"replicas"`
				Template struct {
					Spec struct {
						Containers []struct {
							Resources struct {
								Limits struct {
									Cpu    string `json:"cpu"`
									Memory string `json:"memory"`
								} `json:"limits"`
							} `json:"resources"`
						} `json:"containers"`
					} `json:"spec"`
				} `json:"template"`
			} `json:"spec"`
		} `json:"items"`
	} `json:"data"`
}

func extractDeployResult(data []GenericJson) (res []deployResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid deploy result: %s", err)
	}
	return
}

func cpuNum(cpuStr string) (ret int, err error) {
	if cpuStr[len(cpuStr)-1:] == "m" {
		var cpu int
		cpu, err = strconv.Atoi(cpuStr[:len(cpuStr)-1])
		ret += cpu
	} else {
		var cpu int
		cpu, err = strconv.Atoi(cpuStr)
		ret += 1000 * cpu
	}
	if err != nil {
		err = fmt.Errorf("invalid CPU string")
	}
	return
}

func memNum(memStr string) (ret int, err error) {
	mem, err := strconv.Atoi(memStr[:len(memStr)-2])
	if memStr[len(memStr)-2:] == "Gi" {
		ret += 1024 * mem
	} else {
		ret += mem
	}
	if err != nil {
		err = fmt.Errorf("invalid memory string")
	}
	return
}

func formatDeployPrettyPrint(data []deployResult) (ppc PrettyPrintConfig, err error) {
	ppc.Columns = []string{"NAME", "PODS", "PODS ACTIVE", "CPU", "RAM", "AGE"}
	for _, v := range data {
		for _, item := range v.Data.Items {
			var cpuTotal, memTotal int
			if item.Spec.Replicas != 0 {
				for _, container := range item.Spec.Template.Spec.Containers {
					var cpu, mem int
					cpu, err = cpuNum(container.Resources.Limits.Cpu)
					if err != nil {
						return
					}
					mem, err = memNum(container.Resources.Limits.Memory)
					cpuTotal += cpu
					memTotal += mem
				}
				cpuTotal *= item.Spec.Replicas
				memTotal *= item.Spec.Replicas
			}
			pods := fmt.Sprintf("%d", item.Spec.Replicas)
			if item.Spec.Replicas == 0 {
				pods = "None"
			}
			row := []string{
				item.Metadata.Name,
				pods,
				fmt.Sprintf("%d", item.Status.AvailableReplicas),
				fmt.Sprintf("%dm", cpuTotal),
				fmt.Sprintf("%dMi", memTotal),
				fmt.Sprintf("%dd", int(time.Now().Sub(item.Metadata.CreatedAt).Hours()/24)),
			}
			ppc.Data = append(ppc.Data, row)
		}
	}
	return
}
