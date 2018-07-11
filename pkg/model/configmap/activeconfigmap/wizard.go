package activeconfigmap

import (
	"strings"

	"os"

	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/str"
)

type Config struct {
	EditName  bool
	ConfigMap *configmap.ConfigMap
}

func (c Config) Wizard() configmap.ConfigMap {
	var config configmap.ConfigMap
	if c.ConfigMap != nil {
		config = (*c.ConfigMap).Copy()
	}
	for exit := false; !exit; {
		var menuItems activekit.MenuItems
		if c.EditName {
			menuItems = menuItems.Append(&activekit.MenuItem{
				Label: "Edit name : " + str.Vector{config.Name, "empty"}.FirstNonEmpty(),
				Action: func() error {
					var name = activekit.Promt("Type name, hit Enter to leave %s: ", str.Vector{config.Name, "empty"}.FirstNonEmpty())
					name = strings.TrimSpace(name)
					if name != "" {
						config.Name = name
					}
					return nil
				},
			})
		}
		var configtemsLen = uint(len(config.Data))
		var configItems = config.Items()
		menuItems = menuItems.Append(activekit.ItemsFromIter(configtemsLen, func(index uint) *activekit.MenuItem {
			var configItem = configItems[index]
			return &activekit.MenuItem{
				Label: "Edit " + configItem.Key(),
				Action: func() error {
					var item, status = componentConfigmapItem(configItem)
					switch status {
					case configmapItemComponentStatus{Replace: true}:
						configItems[index] = item
					case configmapItemComponentStatus{Delete: true}:
						configItems = append(configItems[:index], configItems[index+1:]...)
					case configmapItemComponentStatus{DropChanges: true}:
						// pass
					default:
						panic("[configmap.Wizard] edit item: unreachable state")
					}
					return nil
				},
			}
		})...).Append(activekit.MenuItems{
			{
				Label: "Add item",
				Action: func() error {
					var item, status = componentConfigmapItem(configmap.Item{})
					switch status {
					case configmapItemComponentStatus{Replace: true}:
						configItems = append(configItems, item)
					case configmapItemComponentStatus{Delete: true}, configmapItemComponentStatus{DropChanges: true}:
						// pass
					default:
						panic("[configmap.Wizard] add item: unreachable state")
					}
					return nil
				},
			},
			{
				Label: "Confirm",
				Action: func() error {
					if err := ValidateConfigMap(config); err != nil {
						ferr.Println(err)
					} else {
						exit = true
					}
					return nil
				},
			},
			{
				Label: "Exit",
				Action: func() error {
					os.Exit(1)
					return nil
				},
			},
		}...)
		(&activekit.Menu{
			Title: "Configmap",
			Items: menuItems,
		}).Run()
		config.Data = configmap.ConfigMap{}.Data
		config = config.AddItems(configItems...)
	}
	return config
}
