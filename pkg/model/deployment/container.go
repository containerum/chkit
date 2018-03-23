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
		container.Name,
		fmt.Sprintf("CPU %s", container.Limits.CPU),
		fmt.Sprintf("MEM %s", container.Limits.Memory),
	}, ", ")
}
