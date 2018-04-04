package deplactive

import (
	"fmt"
	"strings"

	"git.containerum.net/ch/kube-client/pkg/model"

	"github.com/containerum/chkit/pkg/util/namegen"

	"github.com/containerum/chkit/pkg/model/container"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrInvalidDeploymentName chkitErrors.Err = "invalid deployment name"
)

func getName(defaultName string) string {
	for {
		name, _ := activeToolkit.AskLine(fmt.Sprintf("Print deployment name (or hit Enter to use %q) > ", defaultName))
		if strings.TrimSpace(name) == "" {
			name = defaultName
		}
		if err := validation.ValidateLabel(name); err != nil {
			fmt.Printf("Invalid name %q. Try again\n", name)
			continue
		}
		return name
	}
}

func getReplicas(defaultReplicas uint) uint {
	for {
		replicasStr, _ := activeToolkit.AskLine(fmt.Sprintf("Print number or replicas (1..15, hit Enter to user %d) > ", defaultReplicas))
		replicas := defaultReplicas
		if strings.TrimSpace(replicasStr) == "" {
			return defaultReplicas
		}
		if _, err := fmt.Sscan(replicasStr, &replicas); err != nil || replicas > 15 {
			fmt.Printf("Expected number 1..15! Try again.\n")
			continue
		}
		return replicas
	}
}

func getContainers() []container.Container {
	return nil
}

func getContainer() container.Container {
	con := container.Container{
		model.Container{
			Name:  namegen.Aster() + "-" + namegen.Color(),
			Image: "unknown (required)",
		},
	}
	fmt.Printf("Ok, the hard part. Let's create a container\n")
	for {
		activeToolkit.Options("What's next?", false,
			fmt.Sprintf("Name : %q", con.Name),
			fmt.Sprintf("Image : %q", con.Image))
	}
}
