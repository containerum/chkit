package container

import (
	"fmt"
	"path"
	"strings"

	"github.com/containerum/chkit/pkg/model/limits"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/ninedraft/boxofstuff/str"
)

func (container Container) Validate() error {
	var errs []error
	if err := container.validateName(); err != nil {
		errs = append(errs, err)
	}
	if err := container.validateLimits(); err != nil {
		errs = append(errs, err)
	}
	if err := container.validateEnvs(); err != nil {
		errs = append(errs, err)
	}
	if err := container.validateVolumeMounts(); err != nil {
		errs = append(errs, err)
	}
	if err := container.validateConfigmaps(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid container %q:\n%v\n",
			container.Name,
			str.FromErrs(errs...).Join("\n"))
	}
	return nil
}

func (container Container) validateEnvs() error {
	var errs = make([]error, 0, len(container.Env))
	for _, env := range container.Env {
		if env.Name == "" || strings.Contains(env.Name, " ") {
			errs = append(errs, fmt.Errorf("invalid env name %q", env.Name))
		}
		if env.Value == "" {
			errs = append(errs, fmt.Errorf("invalud env %q value: empty values are not allowed", env.Name, env.Value))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(" + invalid envs:\n",
			str.FromErrs(errs...).
				Map(str.Prefix(" ++ ")).
				Join("\n"))
	}
	return nil
}

func (container Container) validateName() error {
	return validation.ValidateContainerName(container.Name)
}

func (container Container) validateVolumeMounts() error {
	var errs []error
	for _, vol := range container.VolumeMounts {
		if err := validation.ValidateLabel(vol.Name); err != nil {
			errs = append(errs, fmt.Errorf("invalid name %q: %v", vol.Name, err))
		}
		if !path.IsAbs(vol.MountPath) {
			errs = append(errs, fmt.Errorf("invalid mount path %q: expected absolute path", vol.MountPath))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(" + invalid volume mounts:\n%v", str.FromErrs(errs...).
			Map(str.Prefix(" ++ ")).Join("\n"))
	}
	return nil
}

func (container Container) validateConfigmaps() error {
	var errs []error
	for _, configmap := range container.ConfigMaps {
		if err := validation.ValidateLabel(configmap.Name); err != nil {
			errs = append(errs, fmt.Errorf("invalid name %q: %v", configmap.Name, err))
		}
		if !path.IsAbs(configmap.MountPath) {
			errs = append(errs, fmt.Errorf("invalid mount path %q: expected absolute path", configmap.MountPath))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(" + invalid configmaps:\n%v", str.FromErrs(errs...).
			Map(str.Prefix(" ++ ")).
			Join("\n"))
	}
	return nil
}

func (container Container) validateLimits() error {
	var errs []error
	if !limits.CPULimit.Containing(int(container.Limits.CPU)) {
		errs = append(errs, fmt.Errorf("invalid CPU limit %d: expected %v mCPU",
			container.Limits.CPU, limits.CPULimit))
	}
	if !limits.MemLimit.Containing(int(container.Limits.Memory)) {
		errs = append(errs, fmt.Errorf("invalid memory limit %d: expected %v Mb",
			container.Limits.Memory, limits.MemLimit))
	}
	if len(errs) > 0 {
		return fmt.Errorf("%v", str.FromErrs(errs...).
			Map(str.Prefix(" + ")).
			Join("\n"))
	}
	return nil
}
