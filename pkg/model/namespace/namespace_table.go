package namespace

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (Namespace) TableHeaders() []string {
	return []string{"ID", "Label", "Owner", "CPU", "MEM", "Age"}
}

func (namespace Namespace) TableRows() [][]string {
	return [][]string{{
		namespace.ID,
		namespace.Label,
		namespace.Owner,
		namespace.UsageCPU(),
		namespace.UsageMemory(),
		namespace.Age(),
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
