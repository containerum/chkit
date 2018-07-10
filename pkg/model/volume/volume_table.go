package volume

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Volume{}
)

func (Volume) TableHeaders() []string {
	return []string{"Name", "Access", "Storage, GB", "Age"}
}

func (volume Volume) TableRows() [][]string {
	return [][]string{{
		volume.OwnerAndName(),
		func() string {
			var accessCell = fmt.Sprintf("You can %v\n", volume.Access)
			for _, user := range volume.Users {
				accessCell += fmt.Sprintf("%v\n", user.String())
			}
			return accessCell
		}(),
		fmt.Sprintf("%d", volume.Capacity),
		volume.Age(),
	}}
}

func (volume Volume) RenderTable() string {
	return model.RenderTable(volume)
}
