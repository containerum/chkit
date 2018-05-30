package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	solutionsPath     = "/solutions"
	userSolutionsPath = "/user_solutions"
)

// GetSolutionList -- returns list of public solutions
func (client *Client) GetSolutionList() (model.AvailableSolutionsList, error) {
	var solutionList model.AvailableSolutionsList
	err := client.RestAPI.Get(rest.Rq{
		Result: &solutionList,
		URL: rest.URL{
			Path: solutionsPath,
		},
	})
	return solutionList, err
}

func (client *Client) RunSolution(solution model.UserSolution) (model.RunSolutionResponse, error) {
	var resp model.RunSolutionResponse
	err := client.RestAPI.Post(rest.Rq{
		Result: &resp,
		Body:   solution.Copy(),
		URL: rest.URL{
			Path: userSolutionsPath,
		},
	})
	return resp, err
}
