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
	return []string{"Name", "When", "Replicas", "Containers"}
}

func (depl *Deployment) TableRows() [][]string {
	containers := make([]string, len(depl.Containers))
	for _, container := range depl.Containers {
		containers = append(containers,
			fmt.Sprintf("%s", container.String()))
	}
	return [][]string{{
		depl.Name,
		depl.Status.ColumnWhen(),
		depl.Status.ColumnReplicas(),
		strings.Join(containers, "\n"),
	}}
}
