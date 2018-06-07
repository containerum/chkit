package context

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model/namespace"
)

type Namespace struct {
	ID         string
	Label      string
	OwnerLogin string
}

func (ns Namespace) String() string {
	return fmt.Sprintf("%s (%s)", ns.Label, ns.OwnerLogin)
}

func NamespaceFromModel(ns namespace.Namespace) Namespace {
	return Namespace{
		ID:         ns.ID,
		Label:      ns.Label,
		OwnerLogin: ns.OwnerLogin,
	}
}

func (ns Namespace) IsEmpty() bool {
	return ns == (Namespace{})
}
