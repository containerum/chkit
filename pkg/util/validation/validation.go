package validation

import (
	"regexp"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

const (
	ErrInvalidLabel chkitErrors.Err = "invalid label"
)

var (
	labelRe = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
)

func ValidateLabel(label string) error {
	if !labelRe.MatchString(label) {
		return ErrInvalidLabel
	}
	return nil
}
