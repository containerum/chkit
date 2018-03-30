package namespace

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (_ Namespace) TableHeaders() []string {
	return []string{"Label", "CPU", "MEM", "Age"}
}

func (namespace Namespace) TableRows() [][]string {
	age := ""
	if namespace.CreatedAt != nil {
		age = model.TimestampFormat(*namespace.CreatedAt)
	}
	volumes := ""
	for i, volume := range namespace.Volumes {
		if i > 0 {
			volumes += "\n" + volume.Label
		}
		volumes += volume.Label
	}
	return [][]string{{
		namespace.Label,
		fmt.Sprintf("%s/%s",
			namespace.Resources.Used.CPU,
			namespace.Resources.Hard.CPU),
		fmt.Sprintf("%s/%s",
			namespace.Resources.Used.Memory,
			namespace.Resources.Hard.Memory),
		age,
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
