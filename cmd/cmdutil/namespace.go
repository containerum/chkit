package cmdutil

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrNoNamespaces chkitErrors.Err = "no namespaces in account"
)

// GetNamespace -- return namespace for usage in commands. Priority: "namespace" flag, config file.
func GetNamespace(ctx *cli.Context) string {
	if ctx.IsSet("namespace") {
		return ctx.String("namespace")
	}
	return GetConfig(ctx).DefaultNamespace
}

// GetFirstClientNamespace -- fetches namespace list and returns first element. Needed for login.
func GetFirstClientNamespace(ctx *cli.Context) (string, error) {
	nsList, err := GetClient(ctx).GetNamespaceList()
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
