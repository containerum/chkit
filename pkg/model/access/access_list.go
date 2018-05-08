package access

import (
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
)

type AccessList []Access

func AccessListFromNamespaces(nsList namespace.NamespaceList) AccessList {
	var list = make([]Access, 0, len(nsList))
	for _, ns := range nsList {
		list = append(list, AccessFromNamespace(ns))
	}
	return list
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
