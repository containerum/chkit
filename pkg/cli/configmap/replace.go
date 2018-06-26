package cliconfigmap

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/configmap/activeconfigmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force    bool     `flag:"force, f" desc:"suppress confirmation"`
		Item     []string `desc:"configmap item: $KEY:$VALUE"`
		FileItem []string `desc:"configmap file item: $KEY:$FILENAME"`
		File     string   `desc:"file with configmap data, .json, .yaml, .yml"`
	}
	var cmd = &cobra.Command{
		Use:     "configmap",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var cmName string
			var cmList configmap.ConfigMapList
			switch len(args) {
			case 0:
				var err error
				list, err := ctx.Client.GetConfigmapList(ctx.Namespace.ID)
				if err != nil {
					fmt.Println(err)
					ctx.Exit(1)
				}
				(&activekit.Menu{
					Title: "Select configmap",
					Items: activekit.StringSelector(list.Names(), func(name string) error {
						cmName = name
						return nil
					}),
				}).Run()
				cmList = list
			case 1:
				cmName = args[0]
			default:
				cmd.Help()
				return
			}
			if cmList == nil {
				list, err := ctx.Client.GetConfigmapList(ctx.Namespace.ID)
				if err != nil {
					fmt.Println(err)
					ctx.Exit(1)
				}
				cmList = list
			}
			var err error
			cm, ok := cmList.GetByName(cmName)
			if !ok {
				fmt.Printf("configmap %q not found", cm.Name)
				ctx.Exit(1)
			}
			if flags.FileItem != nil || flags.Item != nil {
				var items = make([]configmap.Item, 0, len(flags.Item)+len(flags.FileItem))
				for _, itemString := range flags.Item {
					var item, err = newItem(itemString)
					if err != nil {
						fmt.Println(err)
						ctx.Exit(1)
					}
					items = append(items, item)
				}
				for _, itemString := range flags.Item {
					var item, err = newItem(itemString)
					if err != nil {
						fmt.Println(err)
						ctx.Exit(1)
					}
					content, err := ioutil.ReadFile(item.Value())
					if err != nil {
						fmt.Printf("error while loading file item %q:\n%v\n", itemString, err)
						ctx.Exit(1)
					}
					items = append(items, item.WithValue(string(content)))
				}
				cm = cm.AddItems(items...)
			}
			if flags.File != "" {
				fmt.Printf("Loading configmap from %q\n", flags.File)
				cm, err = activeconfigmap.FromFile(flags.File)
				if err != nil {
					fmt.Println(err)
					ctx.Exit(1)
				}
			}
			if !flags.Force {
				cm = activeconfigmap.Config{
					EditName:  false,
					ConfigMap: &cm,
				}.Wizard()
			}
			if flags.Force ||
				activekit.YesNo("Do you really want to replace configmap %q on server?", cmName) {
				if err := ctx.Client.ReplaceConfigmap(ctx.Namespace.ID, cm); err != nil {
					fmt.Println(err)
					ctx.Exit(1)
				}
			}
		},
	}
	if err := gpflag.ParseTo(&flags, cmd.Flags()); err != nil {
		panic(err)
	}
	return cmd
}

func newItem(itemsString string) (configmap.Item, error) {
	var tokens = strings.SplitN(itemsString, ":", 2)
	if len(tokens) != 2 {
		return configmap.Item{}, fmt.Errorf("unable to parse %q to configmap item", itemsString)
	}
	return configmap.NewItem(strings.TrimSpace(tokens[0]),
		strings.TrimSpace(tokens[1])), nil
}
