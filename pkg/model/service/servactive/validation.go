package servactive

import (
	"bytes"
	"net/url"

	"github.com/containerum/chkit/pkg/util/validation"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/service"
)

const (
	ErrInvalidServiceName    chkitErrors.Err = "invalid service name"
	ErrInvalidServiceDomain  chkitErrors.Err = "invalid service domain"
	ErrInvalidDeploymentName chkitErrors.Err = "invalid deployment name"
)

type ListErr []error

func (list ListErr) Error() string {
	buf := bytes.NewBuffer(make([]byte, 0, 32*len(list)))
	for _, err := range list {
		if _, er := buf.WriteString(err.Error() + "\n"); er != nil {
			return er.Error()
		}
	}
	return buf.String()
}

func validateService(service service.Service) error {
	var errs ListErr
	if err := validation.ValidateLabel(service.Name); err != nil {
		errs = append(errs, ErrInvalidServiceName)
	}
	if service.Domain != "" {
		if _, err := url.Parse(service.Domain); err != nil {
			errs = append(errs, ErrInvalidServiceDomain.Wrap(err))
		}
	}
	if err := validation.ValidateLabel(service.Deploy); err != nil {
		errs = append(errs, ErrInvalidDeploymentName)
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
