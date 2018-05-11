package activeconfigmap

import (
	"fmt"

	"strings"

	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/interview"
	"github.com/containerum/chkit/pkg/util/text"
)

func itemMenu(item configmap.Item) *configmap.Item {
	var oldItem = item
	var ok = false
	var del = false
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: fmt.Sprintf("Item %s", text.Crop(interview.View([]byte(item.Value)), 64)),
			Items: activekit.MenuItems{
				{
					Label: fmt.Sprintf("Edit name  : %s",
						activekit.OrString(item.Key, "undefined, required")),
					Action: func() error {
						var key = activekit.Promt("Type name (hit Enter to leave %s): ",
							activekit.OrString(item.Key, "empty"))
						key = strings.TrimSpace(key)
						if ok := configmap.KeyRegexp().MatchString(key); key != "" && ok {
							item.Key = key
						} else if !ok {
							fmt.Printf("Invalid key %q: must match %q\n", key, configmap.KeyRegexp())
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Edit value :	 %q", text.Crop(interview.View([]byte(item.Value)), 64)),
					Action: func() error {
						item.Value = itemValueMenu(item.Value)
						return nil
					},
				},
				{
					Label: "Delete",
					Action: func() error {
						if activekit.YesNo("Are you sure?") {
							del = true
							exit = true
							ok = false
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						ok = true
						exit = true
						del = false
						return nil
					},
				},
				{
					Label: "Return to previous menu",
					Action: func() error {
						ok = false
						exit = true
						del = false
						return nil
					},
				},
			},
		}).Run()
	}
	if del {
		return nil
	}
	if ok {
		return &item
	}
	return &oldItem
}
