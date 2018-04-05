package validation

import (
	"regexp"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/docker/distribution/reference"
)

const (
	ErrInvalidLabel         chkitErrors.Err = "invalid label"
	ErrInvalidImageName     chkitErrors.Err = "invalid image name"
	ErrInvalidContainerName chkitErrors.Err = "invalid container name"
)

var (
	labelRe         = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	containerNameRe = regexp.MustCompile("^[a-z0-9](([a-z0-9-[^-])){1,61}[a-z0-9]$")
)

func ValidateContainerName(name string) error {
	if !containerNameRe.MatchString(name) {
		return ErrInvalidContainerName
	}
	return nil
}

func ValidateImageName(image string) error {
	if !reference.NameRegexp.MatchString(image) {
		return ErrInvalidImageName
	}
	return nil
}

func ValidateLabel(label string) error {
	if !labelRe.MatchString(label) {
		return ErrInvalidLabel
	}
	return nil
}
