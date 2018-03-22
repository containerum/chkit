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
	volumes := make([]string, len(depl.Volumes))
	for _, volume := range depl.Volumes {
		volumes = append(volumes,
			fmt.Sprintf("%q %s %dGb",
				volume.Label,
				volume.Access,
				volume.Storage))

	}
	return [][]string{{
		depl.Name,
		depl.Status.String(),
		strings.Join(volumes, "\n"),
	}}
}
