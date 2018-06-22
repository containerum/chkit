package deplactive

import (
	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/limits"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrInvalidContainer  chkitErrors.Err = "invalid container"
	ErrInvalidDeployment chkitErrors.Err = "invalid deployment"
)

func ValidateContainer(cont container.Container) error {
	var errs []error
	if err := validation.ValidateLabel(cont.Name); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid container name: %v", err))
	}

	if err := validation.ValidateImageName(cont.Image); err != nil || cont.Image == "" {
		errs = append(errs, fmt.Errorf("\n + invalid image name: %v", err))
	}

	if !limits.CPULimit.Containing(int(cont.Limits.CPU)) {
		errs = append(errs, fmt.Errorf("\n + invald CPU limit %d: must be in %v mCPU", cont.Limits.CPU, limits.CPULimit))
	}

	if !limits.MemLimit.Containing(int(cont.Limits.Memory)) {
		errs = append(errs, fmt.Errorf("\n + invalid memory limit: must be in %v Mb", limits.MemLimit))
	}

	if len(errs) > 0 {
		return ErrInvalidContainer.CommentF("label=%q", cont.Name).AddReasons(errs...)
	}
	return nil
}

func ValidateDeployment(depl deployment.Deployment) error {
	var errs []error
	if !limits.ReplicasLimit.Containing(depl.Replicas) {
		errs = append(errs, fmt.Errorf("\n + invalid replicas number %d: must be %v", depl.Replicas, limits.ReplicasLimit))
	}
	if len(depl.Containers) == 0 {
		errs = append(errs, fmt.Errorf("\n + can't create deployment without containers!"))
	}
	for _, cont := range depl.Containers {
		if err := ValidateContainer(cont); err != nil {
			errs = append(errs, fmt.Errorf("\n + %s", text.Indent(err.Error(), 2)))
		}
	}
	if len(errs) > 0 {
		return ErrInvalidDeployment.CommentF("label=%q", depl.Name).AddReasons(errs...)
	}
	return nil
}
