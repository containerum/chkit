package image

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrInvalidUpdateImage chkitErrors.Err = "invalid update image"
)

func ValidateImage(img model.UpdateImage) error {
	var errs []error
	if err := validation.ValidateImageName(img.Image); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid image %q", img.Image))
	}
	if err := validation.ValidateLabel(img.Container); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid container label %q", img.Container))
	}
	if len(errs) == 0 {
		return nil
	}
	return ErrInvalidUpdateImage.Wrap(errs...)
}
