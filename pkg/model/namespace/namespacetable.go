package namespace

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (_ Namespace) TableHeaders() []string {
	return []string{"Label", "Created" /* "Access",*/, "Volumes"}
}

func (namespace Namespace) TableRows() [][]string {
	creationTime := ""
	if namespace.CreatedAt != nil {
		creationTime = namespace.CreatedAt.Format(model.CreationTimeFormat)
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
		//namespace.Access,
		volumes,
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
