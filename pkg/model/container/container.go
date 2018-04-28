package container

import (
	"fmt"
	"strings"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Container struct {
	kubeModels.Container
}

func (container Container) String() string {
	return strings.Join([]string{
		container.Name,
		fmt.Sprintf("CPU %d", container.Limits.CPU),
		fmt.Sprintf("MEM %d", container.Limits.Memory),
	}, " ") + "; "
}
