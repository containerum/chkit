package activeingress

import (
	"fmt"

	"net/url"

	"regexp"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/ninedraft/ranger/intranger"
)

const (
	ErrInvalidPath    chkitErrors.Err = "invalid path"
	ErrInvalidRule    chkitErrors.Err = "invalid rule"
	ErrInvalidIngress chkitErrors.Err = "invalid ingress"
)

var HostRe = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")

func ValidateIngress(ingr ingress.Ingress) error {
	var errors []error
	for _, rule := range ingr.Rules {
		if err := ValidateRule(rule); err != nil {
			errors = append(errors, fmt.Errorf("\n + %s", text.Indent(err.Error(), 3)))
		}
	}
	if len(errors) > 0 {
		return ErrInvalidIngress.CommentF("ingress=%q", ingr.Name).AddReasons(errors...)
	}
	return nil
}

func ValidatePath(path ingress.Path) error {
	var errors []error
	portLimits := intranger.IntRanger(1, 65535)
	if err := validation.ValidateLabel(path.ServiceName); err != nil {
		errors = append(errors, fmt.Errorf("\n + invalid service name %q", path.ServiceName))
	}
	if !portLimits.Containing(path.ServicePort) {
		errors = append(errors, fmt.Errorf("\n + invalid port: expect %v, got %d", portLimits, path.ServicePort))
	}
	if _, err := url.Parse(path.Path); err != nil {
		errors = append(errors, fmt.Errorf("\n + invalid path: %v", err))
	}
	if len(errors) > 0 {
		return ErrInvalidPath.CommentF("service=%q", path.ServiceName).AddReasons(errors...)
	}
	return nil
}

func ValidateRule(rule ingress.Rule) error {
	var errors []error
	if !HostRe.MatchString(rule.Host) {
		errors = append(errors, fmt.Errorf("\n + %s", text.Indent("invalid hostname", 3)))
	}
	if len(rule.Paths) < 1 {
		errors = append(errors, fmt.Errorf("\n + %s", text.Indent("no paths provided", 3)))
	}
	for _, path := range rule.Paths {
		if err := ValidatePath(path); err != nil {
			errors = append(errors, fmt.Errorf("\n + %s", text.Indent("+ "+err.Error(), 3)))
		}
	}
	if len(errors) > 0 {
		return ErrInvalidRule.CommentF("host=%q", rule.Host).AddReasons(errors...)
	}
	return nil
}
