package clisetup

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
)

func GetDefaultNS(ctx *context.Context, force bool) error {
	logrus.Debugf("getting user namespaces list")
	list, err := ctx.Client.GetNamespaceList()
	if err != nil {
		logrus.WithError(err).Errorf("unable to get user namespace list")
		fmt.Printf("Unable to get default namespace\n")
		return err
	}
	if len(list) == 0 {
		fmt.Printf("You have no namespaces!\n")
	} else if force {
		ctx.SetNamespace(list[0])
		return nil
	} else {
		var menu []*activekit.MenuItem
		for _, ns := range list {
			menu = append(menu, &activekit.MenuItem{
				Label: ns.LabelAndID(),
				Action: func(ns namespace.Namespace) func() error {
					return func() error {
						ctx.SetNamespace(ns)
						return nil
					}
				}(ns),
			})
		}
		_, err := (&activekit.Menu{
			Title: "Select default namespace",
			Items: menu,
		}).Run()
		return err
	}
	return nil
}
