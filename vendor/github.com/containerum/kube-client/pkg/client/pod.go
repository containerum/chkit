package client

import (
	"github.com/containerum/kube-client/pkg/rest"

	"github.com/containerum/kube-client/pkg/model"
)

const (
	deploymentPodsPath = "/namespaces/{namespace}/deployments/{deployment}/pods"
	podsPath           = "/namespaces/{namespace}/pods"
	podPath            = "/namespaces/{namespace}/pods/{pod}"
)

// DeletePod -- deletes pod in provided namespace
func (client *Client) DeletePod(namespace, pod string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: podPath,
			Params: rest.P{
				"pod":       pod,
				"namespace": namespace,
			},
		},
	})
}

// GetPod -- gets a particular pod by name.
func (client *Client) GetPod(namespace, pod string) (model.Pod, error) {
	var gainedPod model.Pod
	err := client.RestAPI.Get(rest.Rq{
		Result: &gainedPod,
		URL: rest.URL{
			Path: podPath,
			Params: rest.P{
				"namespace": namespace,
				"pod":       pod,
			},
		},
	})
	return gainedPod, err
}

// GetPodList -- returns list of pods in provided namespace
func (client *Client) GetPodList(namespace string) (model.PodsList, error) {
	var podList model.PodsList
	err := client.RestAPI.Get(rest.Rq{
		Result: &podList,
		URL: rest.URL{
			Path: podsPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return podList, err
}

// GetDeploymentPodList -- returns list of pods in provided namespace and deployment
func (client *Client) GetDeploymentPodList(namespace, deployment string) (model.PodsList, error) {
	var podList model.PodsList
	err := client.RestAPI.Get(rest.Rq{
		Result: &podList,
		URL: rest.URL{
			Path: deploymentPodsPath,
			Params: rest.P{
				"namespace":  namespace,
				"deployment": deployment,
			},
		},
	})
	return podList, err
}
