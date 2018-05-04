package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/solution"
)

func (client *Client) GetSolutionList() (solution.SolutionList, error) {
	var gainedList solution.SolutionList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetSolutionList()
		switch {
		case err == nil:
			gainedList = solution.SolutionListFromKube(kubeList)
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			er := client.Auth()
			return true, er
		case cherry.In(err, kubeErrors.ErrAccessError()):
			return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return gainedList, err
}
