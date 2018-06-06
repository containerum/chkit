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
	return []string{"Label", "Version", "Status", "Containers", "Age"}
}

func (depl *Deployment) TableRows() [][]string {
	containers := make([]string, 0, len(depl.Containers))
	for _, container := range depl.Containers {
		containers = append(containers,
			fmt.Sprintf("%s [%s]",
				container.Name,
				container.Image))
	}
	age := "undefined"
	if depl.Status != nil {
		age = model.Age(depl.Status.UpdatedAt)
	}
	return [][]string{{
		depl.Name,
		depl.Version.String(),
		depl.StatusString(),
		strings.Join(containers, "\n"),
		age,
	}}
}
