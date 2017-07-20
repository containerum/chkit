package requestresults

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"chkit-v2/chlib"
)

type singleDeployResult []struct {
	DataType string `json:"DataType"`
	Data     struct {
		chlib.Deploy
	} `json:"data"`
}

type deployListResult []struct {
	Data struct {
		Items []chlib.Deploy `json:"items"`
	} `json:"data"`
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

func (l deployListResult) formatPrettyPrint() (ppc prettyPrintConfig, err error) {
	ppc.Columns = []string{"NAME", "PODS", "PODS ACTIVE", "CPU", "RAM", "AGE"}
	for _, item := range l[0].Data.Items {
		var cpuTotal, memTotal int
		if item.Spec.Replicas != 0 {
			for _, container := range item.Spec.Template.Spec.Containers {
				var cpu, mem int
				cpu, err = cpuNum(container.Resources.Limits.CPU)
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
			ageFormat(time.Now().Sub(*item.Metadata.CreationTimestamp)),
		}
		ppc.Data = append(ppc.Data, row)
	}
	return
}

func (s singleDeployResult) Print() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("deploy field get error")
		}
	}()
	allReplicas := s[0].Data.Spec.Replicas
	status := s[0].Data.Status
	strategy := s[0].Data.Spec.Strategy
	conditions := s[0].Data.Status.Conditions
	containers := s[0].Data.Spec.Template.Spec.Containers

	fmt.Printf("%-30s %s\n", "Name:", s[0].Data.Metadata.Name)
	fmt.Printf("%-30s %s\n", "Namespace:", s[0].Data.Metadata.Namespace)
	fmt.Printf("%-30s %s\n", "CreationtTimeStamp:", s[0].Data.Metadata.CreationTimestamp.Format(time.RFC1123))
	fmt.Println("Labels:")
	for k, v := range s[0].Data.Metadata.Labels {
		fmt.Printf("\t%s=%s\n", k, v)
	}
	fmt.Println("Selectors:")
	for k, v := range s[0].Data.Spec.Selector.MatchLabels {
		fmt.Printf("\t%s=%s\n", k, v)
	}
	replFormat := "%-30s %d %s | %d %s | %d %s | %d %s\n"
	if status.UnavaliableReplicas != 0 {
		fmt.Printf(replFormat, "Replicas:", status.UpdatedReplicas, "updated", status.Replicas,
			"total", allReplicas-status.UnavaliableReplicas, "available", status.UnavaliableReplicas, "unavailable")
	} else {
		fmt.Printf(replFormat, "Replicas:", status.UpdatedReplicas, "updated", status.Replicas,
			"total", status.AvailableReplicas, "available", allReplicas-status.AvailableReplicas, "unavailable")
	}
	fmt.Printf("%-30s %v\n", "Strategy", strategy["type"])
	strategyType := strings.ToLower(strategy["type"].(string)[:1]) + strategy["type"].(string)[1:]
	fmt.Printf("%-30s %v max unavailable, %v max surge\n", strategy["type"].(string)+"Strategy",
		strategy[strategyType].(map[string]interface{})["maxUnavailable"],
		strategy[strategyType].(map[string]interface{})["maxSurge"])
	fmt.Println("Conditions:")
	conditionsTable := prettyPrintConfig{
		Columns: []string{"TYPE", "STATUS", "REASON"},
	}
	for _, v := range conditions {
		row := []string{v.Type, v.Status, v.Reason}
		conditionsTable.Data = append(conditionsTable.Data, row)
	}
	conditionsTable.Print()
	fmt.Println("Containers:")
	for _, c := range containers {
		fmt.Printf("\t%s\n", c.Name)
		if len(c.Command) != 0 {
			fmt.Printf("\t\t%-20s %s\n", "Command:", strings.Join(c.Command, " "))
		}
		fmt.Println("\t\tPorts:")
		if len(c.Ports) != 0 {
			ppc := prettyPrintConfig{
				Columns: []string{"Name", "Protocol", "ContPort"},
			}
			for _, p := range c.Ports {
				row := []string{
					p.Name,
					p.Protocol,
					strconv.Itoa(p.ContainerPort),
				}
				ppc.Data = append(ppc.Data, row)
			}
			ppc.Print()
		}
		fmt.Println("\t\tResourceLimit:")
		fmt.Printf("\t\t\t%-10s %s\n", "CPU:", c.Resources.Limits.CPU)
		fmt.Printf("\t\t\t%-10s %s\n", "Memory:", c.Resources.Limits.Memory)
		fmt.Printf("\t\t%-20s %s\n", "Image:", c.Image)
		fmt.Printf("\t\t%-20s %s\n", "ImagePullPolicy:", c.ImagePullPolicy)
	}
	return
}

func init() {
	resultKinds["Deployment"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res singleDeployResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid deployment response: %s", err)
		}
		return res, nil
	}
	resultKinds["DeploymentList"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res deployListResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid deployment list response: %s", err)
		}
		return res.formatPrettyPrint()
	}
}
