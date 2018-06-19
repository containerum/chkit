package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"git.containerum.net/ch/resource-service/pkg/rsErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
)

const (
	ErrUnableToRunAllSolutionComponents chkitErrors.Err = "unable to run all solution components"
)

func (client *Client) GetSolutionsTemplatesList() (solution.TemplatesList, error) {
	var gainedList solution.TemplatesList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetSolutionsTemplatesList()
		switch {
		case err == nil:
			gainedList = solution.TemplatesListFromKube(kubeList)
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
		logrus.WithError(err).Errorf("unable to get solution templates list")
	}
	return gainedList, err
}

func (client *Client) GetSolutionsTemplatesEnvs(template, branch string) (model.SolutionEnv, error) {
	var kubeList model.SolutionEnv
	err := retry(4, func() (bool, error) {
		var err error
		kubeList, err = client.kubeAPIClient.GetSolutionsTemplateEnv(template, branch)
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
		logrus.WithError(err).Errorf("unable to get solution template envs")
	}
	return kubeList, err
}

func (client *Client) RunSolution(sol solution.Solution) error {
	err := retry(0, func() (bool, error) {
		kubeResp, err := client.kubeAPIClient.RunSolution(sol.ToKube(), sol.Namespace, sol.Branch)
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

func (client *Client) GetRunningSolution(namespace, solutionName string) (solution.Solution, error) {
	var sol solution.Solution
	err := retry(4, func() (bool, error) {
		kubeSol, err := client.kubeAPIClient.GetSolution(namespace, solutionName)
		switch {
		case err == nil:
			sol = solution.SolutionFromKube(kubeSol)
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
		logrus.WithError(err).Errorf("unable to get solution")
	}
	return sol, err
}

func (client *Client) GetRunningSolutionsList(namespace string) (solution.SolutionsList, error) {
	var gainedList solution.SolutionsList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetSolutionsNamespaceList(namespace)
		switch {
		case err == nil:
			gainedList = solution.SolutionsListFromKube(kubeList)
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

func (client *Client) GetSolutionDeployments(namespace, solutionName string) (deployment.DeploymentList, error) {
	var list deployment.DeploymentList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetSolutionDeployments(namespace, solutionName)
		switch {
		case err == nil:
			list = deployment.DeploymentListFromKube(kubeList)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrNamespaceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return list, err
}

func (client *Client) GetSolutionServices(namespace, solutionName string) (service.ServiceList, error) {
	var gainedList service.ServiceList
	err := retry(4, func() (bool, error) {
		kubeLsit, err := client.kubeAPIClient.GetSolutionServices(namespace, solutionName)
		switch {
		case err == nil:
			gainedList = service.ServiceListFromKube(kubeLsit)
			return false, err
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrNamespaceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return gainedList, err
}

func (client *Client) DeleteSolution(namespace, solutionName string) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteSolution(namespace, solutionName)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableDeleteResource()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).
			Errorf("unable to get solution")
	}
	return err
}
