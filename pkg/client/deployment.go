package chClient

import (
	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/coblog"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

const (
	ErrContainerAlreadyExists chkitErrors.Err = "container already exists in deployment"
	ErrContainerDoesNotExist  chkitErrors.Err = "container does not exist"
)

func (client *Client) GetDeployment(namespace, deplName string) (deployment.Deployment, error) {
	var depl deployment.Deployment
	err := retry(4, func() (bool, error) {
		kubeDeployment, err := client.kubeAPIClient.GetDeployment(namespace, deplName)
		if err == nil {
			depl = deployment.DeploymentFromKube(kubeDeployment)
		}
		return HandleErrorRetry(client, err)
	})
	return depl, err
}

func (client *Client) GetDeploymentList(namespace string) (deployment.DeploymentList, error) {
	var list deployment.DeploymentList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetDeploymentList(namespace)
		if err == nil {
			list = deployment.DeploymentListFromKube(kubeList)
		}
		return HandleErrorRetry(client, err)
	})
	return list, err
}

func (client *Client) DeleteDeployment(namespace, deplName string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteDeployment(namespace, deplName)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) CreateDeployment(ns string, depl deployment.Deployment) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.CreateDeployment(ns, depl.ToKube())
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) SetContainerImage(ns, depl string, image kubeModels.UpdateImage) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.SetContainerImage(ns, depl, image)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) ReplaceDeployment(ns string, newDepl deployment.Deployment) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.ReplaceDeployment(ns, newDepl.ToKube())
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) GetDeploymentVersions(namespaceID, deploymentName string) (deployment.DeploymentList, error) {
	var list deployment.DeploymentList
	var logger = coblog.Std.Component("chClient.GetDeploymentVersions")
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetDeploymentVersions(namespaceID, deploymentName)
		if err == nil {
			list = deployment.DeploymentListFromKube(kubeList)
		}
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logger.WithError(err).WithField("namespace", namespaceID).
			Errorf("unable to get versions of deployment %q", deploymentName)
	}
	return list, err
}

func (client *Client) ReplaceDeploymentContainer(ns, deplName string, cont container.Container) error {
	var depl, err = client.GetDeployment(ns, deplName)
	if err != nil {
		return err
	}
	var updated, ok = depl.Containers.Replace(cont)
	if !ok {
		return ErrResourceNotExists.CommentF("container %q not found in deployment %q", cont.Name, depl.Name)
	}
	depl.Containers = updated
	return client.ReplaceDeployment(ns, depl)
}

func (client *Client) CreateDeploymentContainer(ns, deplName string, cont container.Container) error {
	var depl, err = client.GetDeployment(ns, deplName)
	if err != nil {
		return err
	}
	var _, ok = depl.Containers.GetByName(cont.Name)
	if ok {
		return ErrContainerAlreadyExists.CommentF("container:%q, deployment:%q", cont.Name, depl.Name)
	}
	depl.Containers = append(depl.Containers, cont)
	return client.ReplaceDeployment(ns, depl)
}

func (client *Client) DeleteDeploymentContainer(ns, deplName, cont string) error {
	var depl, err = client.GetDeployment(ns, deplName)
	if err != nil {
		return err
	}
	var _, ok = depl.Containers.GetByName(cont)
	if !ok {
		return ErrContainerDoesNotExist.CommentF("container:%q, deployment:%q", cont, depl.Name)
	}
	depl.Containers = depl.Containers.DeleteByName(cont)
	return client.ReplaceDeployment(ns, depl)
}

func (client *Client) GetDeploymentDiffWithPreviousVersion(namespace, deployment string, version semver.Version) (string, error) {
	var diff string
	var err error
	return diff, retry(4, func() (bool, error) {
		diff, err = client.kubeAPIClient.
			GetDeploymentDiffWithPreviousVersion(namespace, deployment, version)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) GetDeploymentDiffBetweenVersions(namespace, deployment string, leftVersion, rightVersion semver.Version) (string, error) {
	var diff string
	var err error
	return diff, retry(4, func() (bool, error) {
		diff, err = client.kubeAPIClient.
			GetDeloymentVersionBetweenVersions(namespace, deployment, leftVersion, rightVersion)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) RunDeploymentVersion(namespace, deployment string, version semver.Version) error {
	return retry(4, func() (bool, error) {
		var err = client.kubeAPIClient.RunDeploymentVersion(namespace, deployment, version)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) DeleteDeploymentVersion(namespace, deployment string, version semver.Version) error {
	return retry(4, func() (bool, error) {
		var err = client.kubeAPIClient.DeleteDeploymentVersion(namespace, deployment, version)
		return HandleErrorRetry(client, err)
	})
}
