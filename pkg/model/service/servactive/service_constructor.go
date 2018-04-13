package servactive

import (
	"fmt"
	"strings"

	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
)

const (
	ErrUserExit             chkitErrors.Err = "user exit"
	ErrUserStoppedSession   chkitErrors.Err = "user stopped session"
	ErrInvalidSymbolInLabel chkitErrors.Err = "invalid symbol in label"
	defaultString                           = "undefined"
)

type ConstructorConfig struct {
	Force       bool
	Deployments []string
	Service     *service.Service
}

func RunInteractveConstructor(config ConstructorConfig) (service.Service, error) {
	var err error
	var serv service.Service
	if config.Service != nil {
		serv = *config.Service
	} else {
		serv = defaultService()
	}
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Name: fmt.Sprintf("Set name  : %s", serv.Name),
					Action: func() error {
						serv.Name = getName(serv.Name)
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set deploy: %s", serv.Deploy),
					Action: func() error {
						deploy, err := getDeploy(config.Deployments)
						if err != nil {
							fmt.Println(err)
							return nil
						}
						serv.Deploy = deploy
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set domain: %s", serv.Domain),
					Action: func() error {
						domain, err := getDomain()
						if err != nil {
							fmt.Println(err)
							return nil
						}
						serv.Domain = domain
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set IPs   : [%s]", strings.Join(serv.IPs, ", ")),
					Action: func() error {
						IPs, err := getIPs()
						if err != nil {
							fmt.Println(err)
							return nil
						}
						serv.IPs = IPs
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set ports : %v", service.PortList(serv.Ports)),
					Action: func() error {
						ports, err := getPorts()
						if err != nil {
							fmt.Println(err)
							return nil
						}
						serv.Ports = ports
						return nil
					},
				},
				{
					Name: "Confirm",
					Action: func() error {
						if err = validateService(serv); err != nil {
							fmt.Printf("Error: %v", err)
							return nil
						}
						exit = true
						return nil
					},
				},
				{
					Name: "Exit",
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

func defaultService() service.Service {
	return service.Service{
		Name:   namegen.ColoredPhysics(),
		Domain: "undefined (optional)",
		IPs:    nil,
		Ports:  nil,
		Deploy: "undefined (required)",
	}
}
