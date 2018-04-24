package servactive

import (
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
)

func portEditorWizard(port service.Port) (service.Port, bool) {
	oldPort := port
	ok := false
	deletePort := false
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				setPortName(&port),
				setPortProto(&port),
				setPortPort(&port),
				setTargetPort(&port),
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
					Label: "Delete port",
					Action: func() error {
						exit = true
						ok = false
						deletePort = true
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
		return port, deletePort
	}
	return oldPort, deletePort
}
