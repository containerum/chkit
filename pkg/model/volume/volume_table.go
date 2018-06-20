package volume

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Volume{}
)

func (Volume) TableHeaders() []string {
	return []string{"Name", "Age", "Storage, GB", "Access"}
}

func (volume Volume) TableRows() [][]string {
	return [][]string{{
		volume.Name,
		volume.Age(),
		fmt.Sprintf("%d", volume.Capacity),
		func() string {
			var accessCell = fmt.Sprintf("You can %v\n", volume.Access)
			for _, user := range volume.Users {
				accessCell += fmt.Sprintf("%v\n", user.String())
			}
			return accessCell
		}(),
	}}
}

func (volume Volume) RenderTable() string {
	return model.RenderTable(volume)
}
