package chlib

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func LoadGenericJsonFromFile(path string) (b []GenericJson, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&b)
	return
}

func GetCmdRequestJson(client *Client, kind, name string) (ret []GenericJson, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("can`t extract field: %s", r)
		}
	}()
	var apiResult TcpApiResult
	switch kind {
	case KindNamespaces:
		apiResult, err = client.Get(KindNamespaces, name, "")
		if err != nil {
			return ret, err
		}
		items := apiResult["results"].([]interface{})
		for _, itemI := range items {
			item := itemI.(map[string]interface{})
			_, hasNs := item["data"].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"]
			if hasNs {
				ret = append(ret, GenericJson(item))
			}
		}
	default:
		apiResult, err := client.Get(kind, name, client.userConfig.Namespace)
		if err != nil {
			return ret, err
		}
		items := apiResult["results"].([]interface{})
		for _, itemI := range items {
			ret = append(ret, itemI.(map[string]interface{}))
		}
	}
	return
}

type PrettyPrintConfig struct {
	Columns []string
	Data    [][]string
}

type NsResult struct {
	Data struct {
		Metadata struct {
			CreatedAt time.Time `json:"creationTimestamp"`
			Namespace string    `json:"namespace"`
		} `json:"metadata"`
		Status struct {
			Hard struct {
				LimitsCpu string `json:"limits.cpu"`
				LimitsMem string `json:"limits.memory"`
			} `json:"hard"`
			Used struct {
				LimitsCpu string `json:"limits.cpu"`
				LimitsMem string `json:"limits.memory"`
			} `json:"used"`
		} `json:"status"`
	} `json:"data"`
}

func ExtractNsResults(data []GenericJson) (res []NsResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid namespace response: %s", err)
	}
	return
}

func FormatNamespacePrettyPrint(data []NsResult) (ppc PrettyPrintConfig) {
	ppc.Columns = []string{"NAME", "HARD CPU", "HARD MEMORY", "USED CPU", "USED MEMORY", "AGE"}
	for _, v := range data {
		row := []string{
			v.Data.Metadata.Namespace,
			v.Data.Status.Hard.LimitsCpu,
			v.Data.Status.Hard.LimitsMem,
			v.Data.Status.Used.LimitsCpu,
			v.Data.Status.Used.LimitsMem,
			fmt.Sprintf("%dd", int(time.Now().Sub(v.Data.Metadata.CreatedAt).Hours()/24)),
		}
		ppc.Data = append(ppc.Data, row)
	}
	return
}

type PodResult struct {
	Data struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Status struct {
				Phase             string    `json:"phase"`
				IP                net.IP    `json:"podIP"`
				StartTime         time.Time `json:"startTime"`
				ContainerStatuses []struct {
					RestartCount int `json:"restartCount"`
				} `json:"containerStatuses"`
			}
		} `json:"items"`
	} `json:"data"`
}

func ExtractPodResults(data []GenericJson) (res []PodResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid pod result: %s", err)
	}
	return
}

func FormatPodPrettyPrint(data []PodResult) (ppc PrettyPrintConfig) {
	ppc.Columns = []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP"}
	for _, v := range data {
		for _, item := range v.Data.Items {
			restarts := 0
			for _, containerStatus := range item.Status.ContainerStatuses {
				restarts += containerStatus.RestartCount
			}
			ipStr := item.Status.IP.String()
			if item.Status.IP == nil {
				ipStr = "None"
			}
			row := []string{
				item.Metadata.Name,
				"-/-",
				item.Status.Phase,
				fmt.Sprintf("%d", restarts),
				fmt.Sprintf("%dd", int(time.Now().Sub(item.Status.StartTime).Hours()/24)),
				ipStr,
			}
			ppc.Data = append(ppc.Data, row)
		}
	}
	return
}

type DeployResult struct {
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

func ExtractDeployResult(data []GenericJson) (res []DeployResult, err error) {
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

func FormatDeployPrettyPrint(data []DeployResult) (ppc PrettyPrintConfig, err error) {
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

type SerivceResult struct {
	Data struct {
		Items []struct {
			Metadata struct {
				Name      string    `json:"name"`
				CreatedAt time.Time `json:"creationTimestamp"`
				Labels    struct {
					IsExternal string `json:"external"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				ClusterIP   net.IP   `json:"clusterIP"`
				DomainHosts []string `json:"domainHosts"`
				Replicas    int      `json:"replicas"`
				Ports       []struct {
					Port       int    `json:"port"`
					TargetPort int    `json:"targetPort"`
					Protocol   string `json:"protocol"`
				} `json:"ports"`
			} `json:"spec"`
			Status struct {
				AvailableReplicas int `json:"availableReplicas"`
			} `json:"status"`
		} `json:"items"`
	} `json:"data"`
}

func ExtractServiceResult(data []GenericJson) (res []SerivceResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid service result: %s", err)
	}
	return
}

func FormatServicePrettyPrint(data []SerivceResult) (ppc PrettyPrintConfig) {
	ppc.Columns = []string{"NAME", "CLUSTER-IP", "EXTERNAL", "HOST", "PORT(S)", "AGE"}
	for _, v := range data {
		for _, item := range v.Data.Items {
			var externalHost string
			if len(item.Spec.DomainHosts) != 0 && item.Metadata.Labels.IsExternal == "true" {
				externalHost = strings.Join(item.Spec.DomainHosts, " ,\n")
			} else {
				externalHost = "--"
			}
			var ports []string
			for _, port := range item.Spec.Ports {
				if port.Port == port.TargetPort {
					ports = append(ports, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
				} else {
					ports = append(ports, fmt.Sprintf("%d:%d/%s", port.Port, port.TargetPort, port.Protocol))
				}
			}
			row := []string{
				item.Metadata.Name,
				item.Spec.ClusterIP.String(),
				item.Metadata.Labels.IsExternal,
				externalHost,
				strings.Join(ports, " ,\n"),
				fmt.Sprintf("%d", int(time.Now().Sub(item.Metadata.CreatedAt).Hours()/24)),
			}
			ppc.Data = append(ppc.Data, row)
		}
	}
	return
}

func PrettyPrint(ppc PrettyPrintConfig, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(ppc.Columns)
	table.AppendBulk(ppc.Data)
	table.Render()
}
