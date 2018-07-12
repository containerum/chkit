package servactive

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/sirupsen/logrus"
)

func ReplaceWizard(config ConstructorConfig) (service.Service, error) {
	var err error
	var serv service.Service
	if config.Service != nil {
		serv = *config.Service
	} else {
		serv = DefaultService()
	}
	if len(config.Deployments) == 1 && serv.Deploy == "" {
		serv.Deploy = config.Deployments[0]
	}
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set deploy: %s",
						activekit.OrString(serv.Deploy, "undefined (required)")),
					Action: func() error {
						deploy := getDeploy(serv.Deploy, config.Deployments)
						serv.Deploy = deploy
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set ports : %v", service.PortList(serv.Ports)),
					Action: func() error {
						ports := editPorts(serv.Ports, config.External)
						serv.Ports = ports
						return nil
					},
				},
				{
					Label: "Print to terminal",
					Action: func() error {
						data, err := serv.RenderYAML()
						if err != nil {
							logrus.WithError(err).Errorf("unable to render service to yaml")
							activekit.Attention(err.Error())
						}
						border := strings.Repeat("_", text.Width(data))
						fmt.Printf("%s\n%s\n%s\n", border, data, border)
						return nil
					},
				},
				{
					Label: "Export service to file",
					Action: func() error {
						var fname = activekit.Promt("Type filename: ")
						fname = strings.TrimSpace(fname)
						if fname != "" {
							if err := (porta.Exporter{OutFile: fname}.Export(serv)); err != nil {
								ferr.Printf("unable to export service:\n%v\n", err)
							}
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						if err = ValidateService(serv); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						exit = true
						return nil
					},
				},
				{
					Label: "Exit",
					Action: func() error {
						os.Exit(0)
						return nil
					},
				},
			},
		}).Run()
	}
	return serv, nil
}
