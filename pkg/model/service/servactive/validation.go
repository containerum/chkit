package servactive

import (
	"regexp"

	"github.com/containerum/chkit/pkg/chkitErrors"

	"github.com/containerum/chkit/pkg/model/service"
)

const (
	ErrInvalidLabel chkitErrors.Err = "invalid label"
)

var (
	labelRe = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
)

func validateService(service service.Service) error {

	return nil
}

func validateLabel(label string) error {
	if !labelRe.MatchString(label) {
		return ErrInvalidLabel
	}
	return nil
}
