package deployment

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = Deployment{}
)

func (depl Deployment) RenderTable() string {
	return model.RenderTable(&depl)
}

func (_ *Deployment) TableHeaders() []string {
	return []string{"Name", "Status", "Volumes"}
}

func (depl *Deployment) TableRows() [][]string {
	containers := make([]string, len(depl.Containers))
	for _, container := range depl.Containers {
		containers = append(containers,
			fmt.Sprintf("%s", container.String()))

	}
	return [][]string{{
		depl.Name,
		depl.Status.String(),
		strings.Join(containers, "\n"),
	}}
}
