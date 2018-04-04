package deplactive

import (
	"fmt"
	"strings"

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
