package volume

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Volume{}
)

func (_ *Volume) TableHeaders() []string {
	return []string{"ID", "Label", "Age", "Access", "Replicas", "Capacity, GB"}
}

func (volume *Volume) TableRows() [][]string {
	return [][]string{{
		volume.ID,
		volume.Label,
		model.Age(volume.CreatedAt),
		volume.Access,
		fmt.Sprintf("%d", volume.Replicas),
		fmt.Sprintf("%d", volume.Capacity),
	}}
}

func (volume *Volume) RenderTable() string {
	return model.RenderTable(volume)
}
