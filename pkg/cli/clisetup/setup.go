package clisetup

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

var Config = struct {
	DebugRequests bool
}{}

const (

	// ErrInvalidUserInfo -- invalid user info"
	ErrInvalidUserInfo chkitErrors.Err = "invalid user info"
	// ErrInvalidAPIurl -- invalid API url
	ErrInvalidAPIurl chkitErrors.Err = "invalid API url"
	// ErrUnableToLoadTokens -- unable to load tokens
	ErrUnableToLoadTokens chkitErrors.Err = "unable to load tokens"
	// ErrUnableToSaveTokens -- unable to save tokens
	ErrUnableToSaveTokens chkitErrors.Err = "unable to save tokens"
)

func SetupAll() error {
	logrus.Debugf("loading config")
	if err := configuration.LoadConfig(); err != nil {
		return err
	}
	logrus.Debugf("setuping config")
	if err := SetupConfig(); err != nil {
		return err
	}
	logrus.Debugf("setuping client")
	if err := SetupClient(); err != nil {
		return err
	}
	logrus.Debugf("API: %q", Context.Client.APIaddr)
	return nil
}
