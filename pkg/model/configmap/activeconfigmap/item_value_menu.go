package activeconfigmap

import (
	"fmt"

	"strings"

	"io/ioutil"

	"os"

	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/eof"
	"github.com/containerum/chkit/pkg/util/interview"
	"github.com/containerum/chkit/pkg/util/text"
)

func itemValueMenu(value string) string {
	var oldValue = value
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: fmt.Sprintf("value : %q", text.Crop(interview.View(value), 64)),
			Items: activekit.MenuItems{
				{
					Label: "Load from file",
					Action: func() error {
						fname := activekit.Promt("Type filename (you can drop changes later, hit Enter to return to previous menu): ")
						fname = strings.TrimSpace(fname)
						if fname != "" {
							data, err := ioutil.ReadFile(fname)
							if err != nil {
								fmt.Println(err)
								return nil
							}
							value = interview.View(data)
						}
						return nil
					},
				},
				{
					Label: "Read from input",
					Action: func() error {
						fmt.Printf("Type or paste data (you can drop changes later, hit %s to end input):\n", eof.COMBO)
						data, err := ioutil.ReadAll(os.Stdin)
						if err != nil {
							fmt.Println(err)
							return nil
						}
						value = string(data)
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						exit = true
						return nil
					},
				},
				{
					Label: "Return to previous menu",
					Action: func() error {
						exit = true
						value = oldValue
						return nil
					},
				},
			},
		}).Run()
	}
	return value
}
