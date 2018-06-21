package validation

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/docker/distribution/reference"
	"github.com/ninedraft/ranger/intranger"
	"github.com/satori/go.uuid"
)

const (
	ErrInvalidLabel         chkitErrors.Err = "invalid label"
	ErrInvalidImageName     chkitErrors.Err = "invalid image name"
	ErrInvalidContainerName chkitErrors.Err = "invalid container name"
	ErrInvalidDNSLabel      chkitErrors.Err = "invalid DNS label: "
)

var (
	dnsLabelRe      = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9\\-]{1,63}[a-zA-Z0-9]$")
	numericRe       = regexp.MustCompile("^[0-9]+$")
	labelRe         = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	containerNameRe = regexp.MustCompile("^[a-z0-9]([a-z0-9-[^-]){1,61}[a-z0-9]$")
)

func ValidateContainerName(name string) error {
	name = strings.TrimSpace(name)
	if !containerNameRe.MatchString(name) {
		return ErrInvalidContainerName
	}
	return nil
}

func ValidateImageName(image string) error {
	image = strings.TrimSpace(image)
	if !reference.NameRegexp.MatchString(image) || image == "" {
		return ErrInvalidImageName
	}
	return nil
}

func ValidateLabel(label string) error {
	if !labelRe.MatchString(label) {
		return fmt.Errorf("%v: must satsify %v", ErrInvalidLabel, labelRe)
	}
	return nil
}

// RFC 952 and RFC 1123
func DNSLabel(label string) error {
	DNSlenLimits := intranger.IntRanger(1, 63)
	if !DNSlenLimits.Containing(len(label)) {
		return ErrInvalidDNSLabel.CommentF("DNS label length can be in range %v", DNSlenLimits)
	}
	if !dnsLabelRe.MatchString(label) {
		return ErrInvalidDNSLabel.Comment(
			"must consist of a-Z 1-9 and '-'(dash) letters",
			"must start and end with a-Z 1-9 letters",
		)
	}
	if numericRe.MatchString(label) {
		return ErrInvalidLabel.CommentF("must not consist of all numeric values")
	}
	return nil
}

func ValidateID(ID string) error {
	_, err := uuid.FromString(ID)
	return err
}
