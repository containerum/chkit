package requestresults

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/olekukonko/tablewriter"
)

type singlePodResult []struct {
	Data struct {
		chlib.Pod
	} `json:"data"`
}

type podListResult []struct {
	Data struct {
		Items []chlib.Pod `json:"items"`
	} `json:"data"`
}

func (p podListResult) formatPrettyPrint() (ppc prettyPrintConfig) {
	ppc.Columns = []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP"}
	for _, item := range p[0].Data.Items {
		restarts := 0
		for _, containerStatus := range item.Status.ContainerStatuses {
			restarts += containerStatus.RestartCount
		}
		ipStr := item.Status.PodIP.String()
		if item.Status.PodIP == nil {
			ipStr = "None"
		}
		row := []string{
			item.Metadata.Name,
			"-/-",
			item.Status.Phase,
			fmt.Sprintf("%d", restarts),
			ageFormat(time.Now().Sub(item.Status.StartTime)),
			ipStr,
		}
		ppc.Data = append(ppc.Data, row)
	}
	ppc.Align = tablewriter.ALIGN_LEFT
	return
}

func (p singlePodResult) Print() (err error) {
	metadata := p[0].Data.Metadata
	containers := p[0].Data.Spec.Containers
	restartPolicy := p[0].Data.Spec.RestartPolicy
	termination := time.Duration(p[0].Data.Spec.TerminationGracePeriodSeconds) * time.Second
	system := p[0].Data.Status
	containerStatuses := p[0].Data.Status.ContainerStatuses
	status := p[0].Data.Status.Conditions

	fmt.Println("Pod:")
	fmt.Printf("\t%-20s %s\n", "UserId:", metadata.CreationTimestamp.Format(time.RFC1123))
	fmt.Println("\tLabel:")
	fmt.Printf("\t\t%-20s %s\n", "App:", metadata.Labels["app"])
	fmt.Printf("\t\t%-20s %s\n", "PodTemplateHash:", metadata.Labels["pod-template-hash"])
	fmt.Printf("\t\t%-20s %s\n", "Role", metadata.Labels["role"])
	fmt.Println("Containers:")
	for _, c := range containers {
		fmt.Printf("\t%s\n", c.Name)
		if len(c.Command) != 0 {
			fmt.Println("\t\t%-20s %s\n", "Command:", strings.Join(c.Command, ""))
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
		if len(c.Env) != 0 {
			fmt.Println("\t\tEnvironment:")
			ppc := prettyPrintConfig{
				Columns: []string{"Name", "Value"},
			}
			for _, e := range c.Env {
				ppc.Data = append(ppc.Data, []string{e.Name, e.Value})
			}
		}
		fmt.Println("\t\tResourceLimit:")
		fmt.Printf("\t\t\t%-10s %s\n", "CPU:", c.Resources.Limits.CPU)
		fmt.Printf("\t\t\t%-10s %s\n", "Memory:", c.Resources.Limits.Memory)
		fmt.Printf("\t\t%-20s %s\n", "Image:", c.Image)
		fmt.Printf("\t\t%-20s %s\n", "ImagePullPolicy:", c.ImagePullPolicy)
		fmt.Println("System:")
		fmt.Printf("\t%-30s %s\n", "PodIP:", system.PodIP)
		fmt.Printf("\t%-30s %s\n", "Phase:", system.Phase)
		if !system.StartTime.Equal(time.Time{}) {
			fmt.Printf("\t%-30s %s\n", "StartTime:", system.StartTime)
		}
		fmt.Printf("\t%-30s %ds\n", "TerminationGracePeriod:", int(termination.Seconds()))
		fmt.Printf("\t%-30s %s\n", "RestartPolicy:", restartPolicy)
		fmt.Println("Container Statuses:")
		if len(containerStatuses) != 0 {
			ppc := prettyPrintConfig{
				Columns: []string{"Name", "Ready", "Restart count"},
			}
			for _, cs := range containerStatuses {
				row := []string{
					cs.Name,
					fmt.Sprint(cs.Ready),
					fmt.Sprint(cs.RestartCount),
				}
				ppc.Data = append(ppc.Data, row)
			}
			ppc.Print()
		}
		fmt.Println("Status:")
		ppc := prettyPrintConfig{
			Columns: []string{"Type", "LastTransitionTime", "Status"},
		}
		for _, s := range status {
			row := []string{
				s.Type,
				s.LastTransitionTime.Format(time.RFC1123),
				s.Status,
			}
			ppc.Data = append(ppc.Data, row)
		}
		ppc.Print()
	}
	return
}

func init() {
	resultKinds["Pod"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res singlePodResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid pod response: %s", err)
		}
		return res, nil
	}
	resultKinds["PodList"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res podListResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid pod list response: %s", err)
		}
		return res.formatPrettyPrint(), nil
	}
}
