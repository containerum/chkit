package access

import (
	permModel "git.containerum.net/ch/permissions/pkg/model"
	"github.com/containerum/chkit/pkg/model"
)

type Access struct {
	User      string      `json:"user"`
	Namespace string      `json:"namespace"`
	Access    AccessLevel `json:"access"`
}

func AccessFromNamespace(namespace permModel.NamespaceWithPermissionsJSON) AccessList {
	var aclist = make([]Access, 0)
	for _, p := range namespace.Permissions {
		aclist = append(aclist, Access{
			User:      p.UserLogin,
			Namespace: namespace.Label,
			Access: func() AccessLevel {
				var lvl, _ = LevelFromString(string(p.CurrentAccessLevel))
				return lvl
			}(),
		})
	}

	return aclist
}

func (access Access) RenderTable() string {
	return model.RenderTable(access)
}

func (Access) TableHeaders() []string {
	return []string{
		"Namespace",
		"Access",
		"User",
	}
}

func (access Access) TableRows() [][]string {
	return [][]string{{
		access.Namespace,
		access.Access.String(),
		access.User,
	}}
}
