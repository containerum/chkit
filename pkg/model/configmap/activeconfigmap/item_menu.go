package activeconfigmap

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/ninedraft/boxofstuff/str"
)

type configmapItemComponentStatus struct {
	Delete      bool
	Replace     bool
	DropChanges bool
}

func componentConfigmapItem(configmapItem configmap.Item) (configmap.Item, configmapItemComponentStatus) {
	var item = configmapItem
	var result configmap.Item
	var status = configmapItemComponentStatus{Replace: true}

	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Configmap -> Items",
			Items: activekit.MenuItems{
				{
					Label: "Edit name : " + str.Vector{item.Key(), "none"}.FirstNonEmpty(),
					Action: func() error {
						var name = activekit.Promt("Type new name, hit Enter to leave %s: ", str.Vector{item.Key(), "empty"}.FirstNonEmpty())
						name = strings.TrimSpace(name)
						switch {
						case name == "":
							// default
						case strings.HasPrefix(name, "$"):
							name = strings.TrimPrefix(name, "$")
							item = configmap.NewItem(name, os.Getenv(name))
						case strings.HasPrefix(name, "file://"):
							name = strings.TrimPrefix(name, "file://")
							var data, err = ioutil.ReadFile(name)
							if err != nil {
								ferr.Println(err)
								os.Exit(1)
								return nil
							}
							item = configmap.NewItem(name, string(data))
						default:
							item = item.WithKey(name)
						}
						return nil
					},
				},
				{
					Label: "Edit value: " + str.Vector{item.Value(), "empty"}.FirstNonEmpty(),
					Action: func() error {
						var value = activekit.Promt("Type value, hit Enter to leave %s: ", str.Vector{text.Crop(item.Value(), 16), "empty"}.FirstNonEmpty())
						value = strings.TrimSpace(value)
						switch {
						case value == "":
							// pass
						case strings.HasPrefix(value, "$"):
							value = strings.TrimPrefix(value, "$")
							item = item.WithValue(os.Getenv(value))
						}
						return nil
					},
				},
				{
					Label: "Load from file",
					Action: func() error {
						var fname = activekit.Promt("Type filename, hit Enter to skip: ")
						fname = strings.TrimSpace(fname)
						if fname != "" {
							var data, err = ioutil.ReadFile(fname)
							if err != nil {
								ferr.Println(err)
								return nil
							}
							item = item.WithValue(string(data))
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						result = item
						exit = true
						status = configmapItemComponentStatus{Replace: true}
						return nil
					},
				},
				{
					Label: "Return to previous menu, drop all changes",
					Action: func() error {
						exit = true
						status = configmapItemComponentStatus{DropChanges: true}
						return nil
					},
				},
				{
					Label: "Delete item",
					Action: func() error {
						if activekit.YesNo("Are you sure?") {
							exit = true
							status = configmapItemComponentStatus{Delete: true}
						}
						return nil
					},
				},
			},
		}).Run()
	}
	return result, status
}
