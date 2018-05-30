package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	deploymentsPath = "/namespaces/{namespace}/deployment"
	deploymentPath  = "/namespaces/{namespace}/deployment/{deployment}"
	imagePath       = "/namespaces/{namespace}/deployment/{deployment}/image"
	replicasPath    = "/namespaces/{namespace}/deployment/{deployment}/replicas"
)

// GetDeployment -- consumes a namespace and a deployment names,
// returns a Deployment data OR uninitialized struct AND an error
func (client *Client) GetDeployment(namespace, deployment string) (model.Deployment, error) {
	var depl model.Deployment
	err := client.RestAPI.Get(rest.Rq{
		Result: &depl,
		URL: rest.URL{
			Path: deploymentPath,
			Params: rest.P{
				"namespace":  namespace,
				"deployment": deployment,
			},
		},
	})
	return depl, err
}

// GetDeploymentList -- consumes a namespace and a deployment names,
// returns a list of Deployments OR nil slice AND an error
func (client *Client) GetDeploymentList(namespace string) (model.DeploymentsList, error) {
	var depls model.DeploymentsList
	err := client.RestAPI.Get(rest.Rq{
		Result: &depls,
		URL: rest.URL{
			Path: deploymentsPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return depls, err
}

// DeleteDeployment -- consumes a namespace, a deployment,
// an user role and an ID
func (client *Client) DeleteDeployment(namespace, deployment string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: deploymentPath,
			Params: rest.P{
				"namespace":  namespace,
				"deployment": deployment,
			},
		},
	})
}

// CreateDeployment -- consumes a namespace, an user ID and a Role,
// returns nil if OK
func (client *Client) CreateDeployment(namespace string, deployment model.Deployment) error {
	return client.RestAPI.Post(rest.Rq{
		Body: deployment,
		URL: rest.URL{
			Path: deploymentsPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// SetContainerImage -- set or changes deployment container image
// Consumes namespace, deployment and container data
func (client *Client) SetContainerImage(namespace, deployment string, updateImage model.UpdateImage) error {
	return client.RestAPI.Put(rest.Rq{
		Body: updateImage,
		URL: rest.URL{
			Path: imagePath,
			Params: rest.P{
				"namespace":  namespace,
				"deployment": deployment,
			},
		},
	})
}

// ReplaceDeployment -- replaces deployment in provided namespace with new one
func (client *Client) ReplaceDeployment(namespace string, deployment model.Deployment) error {
	return client.RestAPI.Put(rest.Rq{
		Body: deployment,
		URL: rest.URL{
			Path: deploymentPath,
			Params: rest.P{
				"namespace":  namespace,
				"deployment": deployment.Name,
			},
		},
	})
}

// SetReplicas -- sets or changes deployment replicas
func (client *Client) SetReplicas(namespace, deployment string, replicas int) error {
	return client.RestAPI.Put(rest.Rq{
		Body: model.UpdateReplicas{
			Replicas: replicas,
		},
		URL: rest.URL{
			Path: replicasPath,
			Params: rest.P{
				"namespace":  namespace,
				"deployment": deployment,
			},
		},
	})
}
