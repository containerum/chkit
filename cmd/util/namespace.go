package util

import (
	. "github.com/containerum/chkit/cmd/context"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"
)

const (
	ErrNoNamespaces chkitErrors.Err = "no namespaces in account"
)

// GetFirstClientNamespace -- fetches namespace list and returns first element. Needed for login.
func GetFirstClientNamespace() (string, error) {
	nsList, err := Context.Client.GetNamespaceList()
	if err != nil {
		return "", err
	}
	if len(nsList) <= 0 {
		return "", ErrNoNamespaces
	}
	selectedNS := nsList[0].Label
	logrus.Debugf("Selected namespace \"%s\"", selectedNS)
	return nsList[0].Label, nil
}
