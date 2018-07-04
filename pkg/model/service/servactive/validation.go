package servactive

import (
	"net/url"

	"github.com/containerum/chkit/pkg/util/validation"

	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/service"
)

const (
	ErrInvalidService chkitErrors.Err = "invalid service"
)

func ValidateService(service service.Service) error {
	var errs []error
	if err := validation.ValidateLabel(service.Name); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid service name %q", service.Name))
	}
	if service.Domain != "" {
		if _, err := url.Parse(service.Domain); err != nil {
			errs = append(errs, fmt.Errorf("\n + invalid service domain %q: %v", service.Domain, err))
		}
	}
	if err := validation.ValidateLabel(service.Deploy); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid service deploy %q", service.Deploy))
	}
	if len(service.Ports) == 0 {
		errs = append(errs, fmt.Errorf("\n + none ports found"))
	}
	for _, port := range service.Ports {
		if err := ValidatePort(port); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return ErrInvalidService.Wrap(errs...)
}

func ValidatePort(p service.Port) error {
	var errs []error
	if err := validation.ValidateLabel(p.Name); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid port name %q", p.Name))
	}
	if p.Protocol != "TCP" && p.Protocol != "UDP" {
		errs = append(errs, fmt.Errorf("invalid port protocol %q", p.Protocol))
	}
	if p.TargetPort < 1 || p.TargetPort > 65553 {
		errs = append(errs, fmt.Errorf("invalid target port %d: msut be 1..65553", p.TargetPort))
	}
	if p.Port != nil && (*p.Port < 1 || *p.Port > 65553) {
		errs = append(errs, fmt.Errorf("invalid port %d: must 1..65553", *p.Port))
	}
	if len(errs) == 0 {
		return nil
	}
	return ErrInvalidPort.Wrap(errs...)
}
