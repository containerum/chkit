package chClient

import (
	authErrors "git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	permErrors "git.containerum.net/ch/permissions/pkg/errors"
	rsErrors "git.containerum.net/ch/resource-service/pkg/rsErrors"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/user-manager/pkg/umErrors"
	vmErrors "git.containerum.net/ch/volume-manager/pkg/errors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/chkitErrors"
)

const (
	ErrYouDoNotHaveAccessToResource chkitErrors.Err = "you don't have access to resource"
	ErrResourceNotExists            chkitErrors.Err = "resource not exists"
	ErrFatalError                   chkitErrors.Err = "fatal error"
)

var retriable = []*cherry.Err{
	permErrors.ErrInternal(),
	permErrors.ErrDatabase(),

	umErrors.ErrInternalError(),
	umErrors.ErrLoginFailed(),
	umErrors.ErrLogoutFailed(),

	kubeErrors.ErrInternalError(),
	kubeErrors.ErrUnableGetResourcesList(),
	kubeErrors.ErrUnableGetResource(),
	kubeErrors.ErrUnableCreateResource(),
	kubeErrors.ErrUnableUpdateResource(),
	kubeErrors.ErrUnableDeleteResource(),
	kubeErrors.ErrUnableGetPodLogs(),
	kubeErrors.ErrExecFailure(),
	kubeErrors.ErrVolumeNotReady(),
	kubeErrors.ErrAccessError(),

	rsErrors.ErrInternal(),
	rsErrors.ErrDatabase(),
	rsErrors.ErrUnableCountResources(),

	vmErrors.ErrInternal(),
	vmErrors.ErrDatabase(),

	sErrors.ErrInternalError(),
	sErrors.ErrUnableGetTemplatesList(),
	sErrors.ErrUnableGetSolution(),
	sErrors.ErrUnableCreateSolution(),
	sErrors.ErrUnableDeleteSolution(),
}

var auth = []*cherry.Err{
	authErrors.ErrInvalidToken(),
	authErrors.ErrTokenNotFound(),
	authErrors.ErrTokenNotOwnedBySender(),
}

var notExists = []*cherry.Err{
	kubeErrors.ErrResourceNotExist(),

	rsErrors.ErrResourceNotExists(),

	vmErrors.ErrResourceNotExists(),
}

var noAccess = []*cherry.Err{
	permErrors.ErrResourceNotOwned(),

	rsErrors.ErrPermissionDenied(),
}

func HandleErrorRetry(client *Client, err error) (bool, error) {
	switch {
	case err == nil:
		//No error
		return false, nil
	case cherry.In(err, retriable...):
		//Retriable error
		return true, ErrFatalError.Wrap(err)
	case cherry.In(err, auth...):
		//Auth errors
		return true, client.Auth()
	case cherry.In(err, noAccess...):
		//Resource not exists errors
		return false, ErrResourceNotExists.Wrap(err)
	case cherry.In(err, noAccess...):
		//Resource access errors
		return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
	default:
		//Another error
		return false, ErrFatalError.Wrap(err)
	}
}
