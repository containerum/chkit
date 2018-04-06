package chClient

import "github.com/containerum/chkit/pkg/chkitErrors"

const (
	// ErrYouDoNotHaveAccessToNamespace -- you don't have access to namespace
	ErrYouDoNotHaveAccessToResource chkitErrors.Err = "you don't have access to resource"
	ErrResourceNotExists            chkitErrors.Err = "resource not exists"
	ErrFatalError                   chkitErrors.Err = "fatal error"
)
