package activeconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
)

type Config struct {
	EditName  bool
	ConfigMap *configmap.ConfigMap
}

func (c Config) Wizard() configmap.ConfigMap {
	var config = configmap.ConfigMap{
		Name: namegen.Aster() + "-" + namegen.Physicist(),
	}
	if c.ConfigMap != nil {
		config = *c.ConfigMap
	}
	for exit := false; !exit; {
		var menu activekit.MenuItems
		for _, item := range config.Items() {
			menu = menu.Append(&activekit.MenuItem{
				Label: fmt.Sprintf("Edit %v", item),
				Action: func(item configmap.Item) func() error {
					if i := itemMenu(item); i != nil {
						config.Data[i.Key] = i.Value
					}
					return nil
				}(item),
			})
		}
		(&activekit.Menu{
			Items: func() activekit.MenuItems {
				if c.EditName {
					return activekit.MenuItems{{
						Label: fmt.Sprintf("Edit name : %s",
							activekit.OrString(config.Name, "undefined, required")),
					}}
				}
				return make(activekit.MenuItems, 0, menu.Len())
			}().Append(menu...),
		}).Run()
	}
	return config
}
