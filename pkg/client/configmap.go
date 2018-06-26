package chClient

import (
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/sirupsen/logrus"
)

func (client *Client) CreateConfigMap(ns string, config configmap.ConfigMap) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.CreateConfigMap(ns, config.Name, config.Data)
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to create configmap")
	}
	return err
}

func (client *Client) GetConfigmap(namespace, cmName string) (configmap.ConfigMap, error) {
	var gainedCM configmap.ConfigMap
	err := retry(4, func() (bool, error) {
		kubeConfigmap, err := client.kubeAPIClient.GetConfigMap(namespace, cmName)
		if err == nil {
			gainedCM = configmap.ConfigMapFromKube(kubeConfigmap)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedCM, err
}

func (client *Client) GetConfigmapList(namespace string) (configmap.ConfigMapList, error) {
	var gainedCM configmap.ConfigMapList
	err := retry(4, func() (bool, error) {
		kubeConfigmapList, err := client.kubeAPIClient.GetConfigMapList(namespace)
		if err == nil {
			gainedCM = configmap.ConfigMapListFromKube(kubeConfigmapList)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedCM, err
}

func (client *Client) DeleteConfigmap(namespace, cm string) error {
	var err = retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteConfigMap(namespace, cm)
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to delete configmap %q in %q", cm, namespace)
	}
	return err
}

func (client *Client) ReplaceConfigmap(namespaceID string, cm configmap.ConfigMap) error {
	var err = retry(4, func() (bool, error) {
		err := client.kubeAPIClient.UpdateConfigMap(namespaceID, cm.Name, cm.Copy().Data)
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithField("method", "ReplaceConfigmap").
			WithError(err).
			Errorf("unable to update configmap %q in %q", cm, namespaceID)
	}
	return err
}
