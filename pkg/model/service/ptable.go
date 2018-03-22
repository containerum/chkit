package service

import (
	"strconv"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = new(Port)
)

func (port *Port) RenderTable() string {
	return model.RenderTable(port)
}

func (_ *Port) TableHeaders() []string {
	return []string{"Protocol", "Target port", "Port"}
}

func (port *Port) TableRows() [][]string {
	optionalPort := "none"
	if port.Port != nil {
		optionalPort = strconv.Itoa(*port.Port)
	}
	targetPort := strconv.Itoa(port.TargetPort)
	return [][]string{{
		port.Protocol,
		targetPort,
		optionalPort,
	}}
}
