package chlib

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
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
			Namespace string    `json:"namespace,omitempty"`
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

func PrettyPrint(ppc PrettyPrintConfig, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(ppc.Columns)
	table.AppendBulk(ppc.Data)
	table.Render()
}
