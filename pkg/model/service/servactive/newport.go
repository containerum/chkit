package servactive

import (
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
)

func portCreationWizard(port service.Port, external bool) service.Port {
	oldPort := port
	ok := false
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				setPortName(&port),
				setPortProto(&port),
				setTargetPort(&port),
				setPortPort(&port, external),
				{
					Label: "Confirm",
					Action: func() error {
						if err := validatePort(port); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						ok = true
						exit = true
						return nil
					},
				},
				{
					Label: "Return to previous menu",
					Action: func() error {
						ok = false
						exit = true
						return nil
					},
				},
			},
		}).Run()
	}
	if ok {
		return port
	}
	return oldPort
}
