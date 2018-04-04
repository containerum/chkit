package servactive

import (
	"bytes"
	"net/url"
	"regexp"

	"github.com/containerum/chkit/pkg/chkitErrors"

	"github.com/containerum/chkit/pkg/model/service"
)

const (
	ErrInvalidLabel          chkitErrors.Err = "invalid label"
	ErrInvalidServiceName    chkitErrors.Err = "invalid service name"
	ErrInvalidServiceDomain  chkitErrors.Err = "invalid service domain"
	ErrInvalidDeploymentName chkitErrors.Err = "invalid deployment name"
)

var (
	labelRe = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
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
	if err := validateLabel(service.Name); err != nil {
		errs = append(errs, ErrInvalidServiceName)
	}
	if service.Domain != "" {
		if _, err := url.Parse(service.Domain); err != nil {
			errs = append(errs, ErrInvalidServiceDomain.Wrap(err))
		}
	}
	if err := validateLabel(service.Deploy); err != nil {
		errs = append(errs, ErrInvalidDeploymentName)
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateLabel(label string) error {
	if !labelRe.MatchString(label) {
		return ErrInvalidLabel
	}
	return nil
}
