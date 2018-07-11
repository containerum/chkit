package container

import (
	"fmt"

	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

type Ports []model.ContainerPort

func (ports Ports) Len() int {
	return len(ports)
}

func (ports Ports) Strings() str.Vector {
	var str = make([]string, 0, ports.Len())
	for _, port := range ports {
		str = append(str, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
	}
	return str
}

func (ports Ports) Ports() []int {
	var p = make([]int, 0, ports.Len())
	for _, port := range ports {
		p = append(p, port.Port)
	}
	return p
}

func (ports Ports) TCP() Ports {
	return ports.Filter(func(port model.ContainerPort) bool {
		return port.Protocol == model.TCP
	})
}

func (ports Ports) UDP() Ports {
	return ports.Filter(func(port model.ContainerPort) bool {
		return port.Protocol == model.UDP
	})
}

func (ports Ports) Filter(pred func(port model.ContainerPort) bool) Ports {
	var filtered = make(Ports, 0, ports.Len())
	for _, p := range ports {
		if pred(p) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

type PortMap map[string]model.ContainerPort

func (pmap PortMap) Copy() PortMap {
	var cp = make(PortMap, pmap.Len())
	for _, port := range pmap {
		cp[port.Name] = port
	}
	return cp
}

func (pmap PortMap) Len() int {
	return len(pmap)
}

func NewPortSet(ports Ports) PortMap {
	var pmap = make(PortMap, ports.Len())
	for _, port := range ports {
		pmap[port.Name] = port
	}
	return pmap
}

func (set PortMap) Ports() Ports {
	var ports = make(Ports, 0, len(set))
	for _, port := range set {
		ports = append(ports, port)
	}
	return ports
}

func (pmap PortMap) Merge(rights ...model.ContainerPort) PortMap {
	var left = pmap.Copy()
	for _, right := range rights {
		left[right.Name] = right
	}
	return left
}
