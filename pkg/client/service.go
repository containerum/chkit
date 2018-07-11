package chClient

import (
	"github.com/containerum/chkit/pkg/model/service"
)

func (client *Client) GetService(namespace, serviceName string) (service.Service, error) {
	var gainedService service.Service
	err := retry(4, func() (bool, error) {
		kubeService, err := client.kubeAPIClient.GetService(namespace, serviceName)
		if err == nil {
			gainedService = service.ServiceFromKube(kubeService)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedService, err
}

func (client *Client) GetServiceList(namespace string) (service.ServiceList, error) {
	var gainedList service.ServiceList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetServiceList(namespace)
		if err == nil {
			gainedList = service.ServiceListFromKube(kubeList)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedList, err
}

func (client *Client) DeleteService(namespace, service string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteService(namespace, service)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) CreateService(ns string, serv service.Service) error {
	return retry(4, func() (bool, error) {
		_, err := client.kubeAPIClient.CreateService(ns, serv.ToKube())
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) ReplaceService(ns string, serv service.Service) error {
	return retry(4, func() (bool, error) {
		_, err := client.kubeAPIClient.UpdateService(ns, serv.ToKube())
		return HandleErrorRetry(client, err)
	})
}
