package activeconfigmap

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
	kubeModels "github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
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
	if config.Data == nil {
		config.Data = make(kubeModels.ConfigMapData, 16)
	}
	for exit := false; !exit; {
		var menu activekit.MenuItems
		for _, item := range config.Items() {
			menu = menu.Append(&activekit.MenuItem{
				Label: fmt.Sprintf("Edit %s", text.Crop(item.String(), 64)),
				Action: func(item configmap.Item) func() error {
					return func() error {
						if i := itemMenu(item); i != nil {
							var key, value = i.Data()
							config.Data[key] = base64.StdEncoding.EncodeToString([]byte(value))
						}
						return nil
					}
				}(item),
			})
		}
		(&activekit.Menu{
			Items: func() activekit.MenuItems {
				if c.EditName {
					return activekit.MenuItems{{
						Label: fmt.Sprintf("Edit name : %s",
							activekit.OrString(config.Name, "undefined, required")),
						Action: func() error {
							name := activekit.Promt("Type configmap name (hit Enter to leave %s): ",
								activekit.OrString(config.Name, "empty"))
							name = strings.TrimSpace(name)
							if err := validation.ValidateLabel(name); name != "" && err == nil {
								config.Name = name
							} else if err != nil {
								ferr.Printf("Invalid name %q!\n", name)
							}
							return nil
						},
					}}
				}
				return make(activekit.MenuItems, 0, menu.Len())
			}().Append(
				menu.Append(activekit.MenuItems{
					{
						Label: "Add item",
						Action: func() error {
							if i := itemMenu(configmap.Item{}); i != nil {
								config.Data[i.Key()] = base64.StdEncoding.EncodeToString([]byte(i.ValueDecoded()))
							}
							return nil
						},
					},
					{
						Label: "Print to terminal",
						Action: func() error {
							data, err := config.RenderYAML()
							if err != nil {
								logrus.WithError(err).Errorf("unable to render configmap to yaml")
								activekit.Attention(err.Error())
							}
							border := strings.Repeat("_", text.Width(data))
							fmt.Printf("%s\n%s\n%s\n", border, data, border)
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
				}...)...,
			)}).Run()
	}
	return config
}
