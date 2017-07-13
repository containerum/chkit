package chlib

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type serivceResult struct {
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

func extractServiceResult(data []GenericJson) (res []serivceResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid service result: %s", err)
	}
	return
}

func formatServicePrettyPrint(data []serivceResult) (ppc PrettyPrintConfig) {
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
