package namespace

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (_ Namespace) TableHeaders() []string {
	return []string{"Label", "Created", "Limits"}
}

func (namespace Namespace) TableRows() [][]string {
	creationTime := ""
	if namespace.CreatedAt != nil {
		creationTime = model.TimestampFormat(*namespace.CreatedAt)
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
		creationTime,
		fmt.Sprintf("CPU: %s; MEM %s",
			namespace.Resources.Hard.CPU,
			namespace.Resources.Hard.Memory),
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
