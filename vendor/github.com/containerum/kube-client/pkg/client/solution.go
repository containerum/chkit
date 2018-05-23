package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	solutionListPath  = "/solutions"
	userSolutionsPath = "/user_solutions"
)

// GetSolutionList -- returns list of public solutions
func (client *Client) GetSolutionList() (model.AvailableSolutionsList, error) {
	var solutionList model.AvailableSolutionsList
	err := client.RestAPI.Get(rest.Rq{
		Result: &solutionList,
		URL: rest.URL{
			Path: solutionListPath,
		},
	})
	return solutionList, err
}

func (client *Client) RunSolution(solution model.UserSolution) (model.RunSolutionResponce, error) {
	var resp model.RunSolutionResponce
	err := client.RestAPI.Post(rest.Rq{
		Result: &resp,
		Body:   solution.Copy(),
		URL: rest.URL{
			Path: userSolutionsPath,
		},
	})
	return resp, err
}
