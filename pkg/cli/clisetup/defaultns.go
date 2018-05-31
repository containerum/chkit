package clisetup

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
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
		ctx.Namespace = list[0].ID
		ctx.Changed = true
		return nil
	} else {
		var menu []*activekit.MenuItem
		for _, ns := range list {
			menu = append(menu, &activekit.MenuItem{
				Label: fmt.Sprintf("%s %s", ns.Label, ns.ID),
				Action: func(ns string) func() error {
					return func() error {
						ctx.Namespace = ns
						ctx.Changed = true
						return nil
					}
				}(ns.ID),
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
