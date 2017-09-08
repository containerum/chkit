package requestresults

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/containerum/chkit.v2/chlib"
)

type singleServiceResult []struct {
	Data struct {
		chlib.Service
	} `json:"data"`
}

type serviceListResult []struct {
	Data struct {
		Items []chlib.Service `json:"items"`
	} `json:"data"`
}

var ServiceColumns = []string{"NAME", "CLUSTER-IP", "EXTERNAL", "HOST", "PORTS", "AGE"}

func (s serviceListResult) formatPrettyPrint() (ppc prettyPrintConfig) {
	ppc.Columns = ServiceColumns
	for _, item := range s[0].Data.Items {
		var externalHost string
		external := item.Metadata.Labels["external"]
		if len(item.Spec.DomainHosts) != 0 && external == "true" {
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
			external,
			externalHost,
			strings.Join(ports, " ,\n"),
			ageFormat(time.Now().Sub(*item.Metadata.CreationTimestamp)),
		}
		ppc.Data = append(ppc.Data, row)
	}
	return
}

func (s singleServiceResult) Print() error {
	metadata := s[0].Data.Metadata
	spec := s[0].Data.Spec
	fmt.Printf("%-30s %s\n", "Name:", metadata.Name)
	fmt.Printf("%-30s %s\n", "Namespace:", metadata.Namespace)
	if len(metadata.Labels) != 0 {
		fmt.Println("Labels:")
		for k, v := range metadata.Labels {
			fmt.Printf("\t%s=%s\n", k, v)
		}
	}
	if len(spec.Selector) != 0 {
		fmt.Println("Selectors:")
		for k, v := range spec.Selector {
			fmt.Printf("\t%s=%s\n", k, v)
		}
	}
	fmt.Printf("%-30s %s \n", "Type:", spec.Type)
	fmt.Printf("%-30s %s \n", "IP:", spec.ClusterIP)
	for _, p := range spec.Ports {
		if p.TargetPort == p.Port {
			fmt.Printf("%-30s %d/%d\n", "Port:", p.Port, p.Protocol)
		} else {
			fmt.Printf("%-30s %d:%d/%s\n", "Port:", p.Port, p.TargetPort, p.Protocol)
		}
	}
	isExternal := metadata.Labels["external"]
	fmt.Printf("%-30s %s\n", "External:", isExternal)
	if len(spec.DomainHosts) != 0 && isExternal == "true" {
		fmt.Printf("%-30s %s \n", "External hosts:", strings.Join(spec.DomainHosts, " ,"))
	} else {
		fmt.Printf("%-30s --\n", "External hosts:")
	}
	return nil
}

func init() {
	resultKinds["Service"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res singleServiceResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid pod response: %s", err)
		}
		return res, nil
	}
	resultKinds["ServiceList"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res serviceListResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid pod list response: %s", err)
		}
		return res.formatPrettyPrint(), nil
	}
}
