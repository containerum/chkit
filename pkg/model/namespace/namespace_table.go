package namespace

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (Namespace) TableHeaders() []string {
	return []string{"Label", "ID", "Owner", "CPU", "MEM", "Age"}
}

func (namespace Namespace) TableRows() [][]string {
	return [][]string{{
		namespace.Label,
		namespace.ID,
		namespace.OwnerLogin,
		namespace.UsageCPU(),
		namespace.UsageMemory(),
		namespace.Age(),
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
