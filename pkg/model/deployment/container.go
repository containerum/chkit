package deployment

import (
	"fmt"
	"strings"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
)

type Container struct {
	kubeModels.Container
}

func (container *Container) String() string {
	return strings.Join([]string{
		"Name: " + container.Name,
		"Resources: " + fmt.Sprintf("CPU %s MEM %s",
			container.Limits.CPU,
			container.Limits.Memory),
	}, "\n")
}
