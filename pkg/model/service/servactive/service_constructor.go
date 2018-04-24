package servactive

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
)

type ConstructorConfig struct {
	Force       bool
	Deployments []string
	Service     *service.Service
}

func Wizard(config ConstructorConfig) (service.Service, error) {
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
					Label: fmt.Sprintf("Set name  : %s",
						activekit.OrString(serv.Name, "undefined (required)")),
					Action: func() error {
						serv.Name = getName(serv.Name)
						return nil
					},
				},
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
						ports := editPorts(serv.Ports)
						serv.Ports = ports
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
						if yes, _ := activekit.Yes("Do you really want to exit?"); yes {
							os.Exit(0)
						}
						return nil
					},
				},
			},
		}).Run()
	}
	return serv, nil
}

func DefaultService() service.Service {
	return service.Service{
		Name:   namegen.ColoredPhysics(),
		Domain: "",
		IPs:    nil,
		Ports:  nil,
		Deploy: "",
	}
}
