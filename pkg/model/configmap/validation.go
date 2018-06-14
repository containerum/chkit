package configmap

import (
	"fmt"

	"time"

	"regexp"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrInvalidConfigMap chkitErrors.Err = "invalid config map"
)

var (
	keyRe = regexp.MustCompile("^[-._a-zA-Z0-9]+$")
)

func KeyRegexp() *regexp.Regexp {
	return keyRe.Copy()
}

func (config ConfigMap) Validate() error {
	var errors []error
	if err := validation.ValidateLabel(config.Name); err != nil {
		errors = append(errors, fmt.Errorf(" + invalid name %q\n", config.Name))
	}
	if _, err := time.Parse(model.TimestampFormat, config.CreatedAt); err != nil && config.CreatedAt != "" {
		errors = append(errors, fmt.Errorf(" + invalid timestamp format: %v\n", err))
	}
	if len(config.Data) == 0 {
		errors = append(errors, fmt.Errorf(" + configmap must contains at least one item\n"))
	}
	for _, item := range config.Items() {
		if !keyRe.MatchString(item.Key) {
			errors = append(errors, fmt.Errorf(" + invalid configmap key %q: key must match %q\n", item.Key, keyRe))
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return ErrInvalidConfigMap.Wrap(errors...)
}
