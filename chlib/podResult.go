package chlib

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type podResult struct {
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

func extractPodResults(data []GenericJson) (res []podResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid pod result: %s", err)
	}
	return
}

func formatPodPrettyPrint(data []podResult) (ppc PrettyPrintConfig) {
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
