package access

import (
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
)

type Access struct {
	Namespace string      `json:"namespace"`
	Access    AccessLevel `json:"access"`
}

func AccessFromNamespace(namespace namespace.Namespace) Access {
	return Access{
		Namespace: namespace.Label,
		Access: func() AccessLevel {
			var lvl, _ = LevelFromString(namespace.Access)
			return lvl
		}(),
	}
}

func (access Access) RenderTable() string {
	return model.RenderTable(access)
}

func (Access) TableHeaders() []string {
	return []string{
		"Namespace",
		"Access",
	}
}

func (access Access) TableRows() [][]string {
	return [][]string{{
		access.Namespace,
		access.Access.String(),
	}}
}
