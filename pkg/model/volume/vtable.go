package volume

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Volume{}
)

func (_ *Volume) TableHeaders() []string {
	return []string{"Label", "Created", "Access", "Replicas", "Storage, GB"}
}

func (volume *Volume) TableRows() [][]string {
	return [][]string{{
		volume.Label,
		model.TimestampFormat(volume.CreatedAt),
		volume.Access,
		fmt.Sprintf("%d", volume.Replicas),
		fmt.Sprintf("%d", volume.Storage),
	}}
}

func (volume *Volume) RenderTable() string {
	return model.RenderTable(volume)
}
