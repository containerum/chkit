package client

import (
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	kubeAPIdeploymentPath  = "/namespaces/{namespace}/deployments/{deployment}"
	kubeAPIdeploymentsPath = "/namespaces/{namespace}/deployments"

	resourceDeploymentRootPath = "/namespace/{namespace}/deployment"
	resourceDeploymentPath     = "/namespace/{namespace}/deployment/{deployment}"
	resourceImagePath          = "/namespace/{namespace}/deployment/{deployment}/image"
	resourceReplicasPath       = "/namespace/{namespace}/deployment/{deployment}/replicas"
)

// GetDeployment -- consumes a namespace and a deployment names,
// returns a Deployment data OR uninitialized struct AND an error
func (client *Client) GetDeployment(namespace, deployment string) (model.Deployment, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace":  namespace,
			"deployment": deployment,
		}).SetResult(model.Deployment{}).
		SetError(cherry.Err{}).
		Get(client.APIurl + kubeAPIdeploymentPath)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return model.Deployment{}, err
	}
	return *resp.Result().(*model.Deployment), nil
}

// GetDeploymentList -- consumes a namespace and a deployment names,
// returns a list of Deployments OR nil slice AND an error
func (client *Client) GetDeploymentList(namespace string) ([]model.Deployment, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetResult([]model.Deployment{}).
		SetError(cherry.Err{}).
		Get(client.APIurl + kubeAPIdeploymentsPath)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.Deployment), nil
}

// DeleteDeployment -- consumes a namespace, a deployment,
// an user role and an ID
func (client *Client) DeleteDeployment(namespace, deployment string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace":  namespace,
			"deployment": deployment,
		}).SetError(cherry.Err{}).
		Delete(client.ResourceAddr + resourceDeploymentPath)
	return MapErrors(resp, err, http.StatusOK)
}

// CreateDeployment -- consumes a namespace, an user ID and a Role,
// returns nil if OK
func (client *Client) CreateDeployment(namespace string, deployment model.Deployment) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(deployment).
		SetError(cherry.Err{}).
		Post(client.ResourceAddr + resourceDeploymentRootPath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted)
}

// SetContainerImage -- set or changes deployment container image
// Consumes namespace, deployment and container data
func (client *Client) SetContainerImage(namespace, deployment string, updateImage model.UpdateImage) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace":  namespace,
			"deployment": deployment,
		}).SetBody(updateImage).
		SetError(cherry.Err{}).
		Put(client.ResourceAddr + resourceImagePath)
	return MapErrors(resp, err,
		http.StatusAccepted,
		http.StatusOK,
		http.StatusNoContent)
}

// ReplaceDeployment -- replaces deployment in provided namespace with new one
func (client *Client) ReplaceDeployment(namespace string, deployment model.Deployment) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace":  namespace,
			"deployment": deployment.Name,
		}).SetBody(deployment).
		SetError(cherry.Err{}).
		Put(client.ResourceAddr + resourceDeploymentPath)
	return MapErrors(resp, err, http.StatusOK)
}

// SetReplicas -- sets or changes deployment replicas
func (client *Client) SetReplicas(namespace, deployment string, replicas int) error {
	resp, err := client.Request.SetPathParams(map[string]string{
		"namespace":  namespace,
		"deployment": deployment,
	}).SetBody(model.UpdateReplicas{
		Replicas: replicas,
	}).SetError(cherry.Err{}).
		Put(client.ResourceAddr + resourceReplicasPath)
	return MapErrors(resp, err,
		http.StatusAccepted,
		http.StatusOK)
}
