package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	servicePath  = "/namespaces/{namespace}/services/{service}"
	servicesPath = "/namespaces/{namespace}/services"
)

// GetService -- consume a namespace id and a service name
// returns a Service OR an uninitialized Service struct AND an error
func (client *Client) GetService(namespace, serviceName string) (model.Service, error) {
	var service model.Service
	err := client.RestAPI.Get(rest.Rq{
		Result: &service,
		URL: rest.URL{
			Path: servicePath,
			Params: rest.P{
				"namespace": namespace,
				"service":   serviceName,
			},
		},
	})
	return service, err
}

// GetServiceList -- consumes a namespace name
// returns a slice of Services OR a nil slice AND an error
func (client *Client) GetServiceList(namespace string) (model.ServicesList, error) {
	var serviceList model.ServicesList
	err := client.RestAPI.Get(rest.Rq{
		Result: &serviceList,
		URL: rest.URL{
			Path: servicesPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return serviceList, err
}

// CreateService -- consumes a namespace name and a Service data,
// returns the created Service AND nil OR an uninitialized Service AND an error
func (client *Client) CreateService(namespace string, service model.Service) (model.Service, error) {
	var gainedService model.Service
	err := client.RestAPI.Post(rest.Rq{
		Body:   service,
		Result: &gainedService,
		URL: rest.URL{
			Path: servicesPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return gainedService, err
}

// DeleteService -- consumes a namespace, a servicce name,
// returns error in case of problem
func (client *Client) DeleteService(namespace, serviceName string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: servicePath,
			Params: rest.P{
				"namespace": namespace,
				"service":   serviceName,
			},
		},
	})
}

// UpdateService -- consumes a namespace, a service data,
// returns an ipdated Service OR an uninitialized Service AND an error
func (client *Client) UpdateService(namespace string, service model.Service) (model.Service, error) {
	var gainedService model.Service
	err := client.RestAPI.Put(rest.Rq{
		Body:   service,
		Result: &gainedService,
		URL: rest.URL{
			Path: servicePath,
			Params: rest.P{
				"namespace": namespace,
				"service":   service.Name,
			},
		},
	})
	return gainedService, err
}
