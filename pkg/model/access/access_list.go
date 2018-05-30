package access

import (
	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type AccessList []Access

func (list AccessList) Len() int {
	return len(list)
}

func (list AccessList) New() AccessList {
	return make(AccessList, 0, list.Len())
}

func (list AccessList) Copy() AccessList {
	return append(list.New(), list...)
}

func (list AccessList) RenderTable() string {
	return model.RenderTable(list)
}

func (AccessList) TableHeaders() []string {
	return Access{}.TableHeaders()
}

func (list AccessList) TableRows() [][]string {
	var rows = make([][]string, 0, len(list))
	for _, access := range list {
		rows = append(rows, access.TableRows()...)
	}
	return rows
}

func (list AccessList) Filter(pred func(Access) bool) AccessList {
	var filtered = list.New()
	for _, access := range list {
		if pred(access) {
			filtered = append(filtered, access)
		}
	}
	return filtered
}

func (list AccessList) Names() []string {
	var names = make([]string, 0, list.Len())
	for _, access := range list {
		names = append(names, access.Username)
	}
	return names
}

func (list AccessList) GetByLevels(levels ...kubeModels.AccessLevel) AccessList {
	return list.Filter(func(access Access) bool {
		for _, lvl := range levels {
			if lvl == access.AccessLevel {
				return true
			}
		}
		return false
	})
}

func (list AccessList) GetOwner() string {
	return list.Filter(func(access Access) bool {
		return access.AccessLevel == kubeModels.Owner
	})[0].Username
}
