package namespace

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (Namespace) TableHeaders() []string {
	return []string{"Label", "Access level", "ID", "CPU", "MEM", "Age"}
}

func (namespace Namespace) TableRows() [][]string {
	return [][]string{{
		namespace.OwnerAndLabel(),
		namespace.Access.String(),
		namespace.ID,
		namespace.UsageCPU(),
		namespace.UsageMemory(),
		namespace.Age(),
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
