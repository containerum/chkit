package activesolution

import (
	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrInvalidSolution chkitErrors.Err = "invalid ingress"
)

func ValidateIngress(sol solution.Solution) error {
	var errs []error
	if err := validation.ValidateLabel(sol.Name); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid solution name %q", sol.Name))
	}
	if sol.Template == "" {
		errs = append(errs, fmt.Errorf("\n + invalid solution template %q", sol.Template))
	}
	if len(errs) > 0 {
		return ErrInvalidSolution.CommentF("name=%q", sol.Name).AddReasons(errs...)
	}
	return nil
}
