package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
)

const (
	ErrUnableToRunAllSolutionComponents chkitErrors.Err = "unable to run all solution components"
)

func (client *Client) GetSolutionsTemplatesList() (solution.SolutionList, error) {
	var gainedList solution.SolutionList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetSolutionsTemplatesList()
		switch {
		case err == nil:
			gainedList = solution.SolutionListFromKube(kubeList)
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			logrus.Debugf("running auth")
			er := client.Auth()
			return true, er
		case cherry.In(err, kubeErrors.ErrAccessError()):
			return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to get solution list")
	}
	return gainedList, err
}

func (client *Client) GetSolutionsTemplatesEnvs(template string) (model.SolutionEnv, error) {
	var kubeList model.SolutionEnv
	err := retry(4, func() (bool, error) {
		var err error
		kubeList, err = client.kubeAPIClient.GetSolutionsTemplateEnv(template)
		switch {
		case err == nil:
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			logrus.Debugf("running auth")
			er := client.Auth()
			return true, er
		case cherry.In(err, kubeErrors.ErrAccessError()):
			return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to get solution list")
	}
	return kubeList, err
}

func (client *Client) RunSolution(sol solution.UserSolution) error {
	err := retry(4, func() (bool, error) {
		kubeResp, err := client.kubeAPIClient.RunSolution(sol.ToKube(), sol.Namespace)
		switch {
		case err == nil:
			if kubeResp.NotCreated != 0 || len(kubeResp.Errors) != 0 {
				return false, ErrUnableToRunAllSolutionComponents.Comment(kubeResp.Errors...)
			}
			return false, nil
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			er := client.Auth()
			logrus.Debugf("running auth")
			return true, er
		case cherry.In(err, kubeErrors.ErrAccessError()):
			return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to run solution")
	}
	return err
}
