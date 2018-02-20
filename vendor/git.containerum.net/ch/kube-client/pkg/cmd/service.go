package cmd

import (
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	servicePath  = "/namespaces/{namespace}/services/{service}"
	servicesPath = "/namespaces/{namespace}/services"
)

// GetService -- consume a namespace id and a service name
// returns a Service OR an uninitialized Service struct AND an error
func (client *Client) GetService(namespace, serviceName string) (model.Service, error) {
	resp, err := client.Request.
		SetResult(model.Service{}).
		SetPathParams(map[string]string{
			"namespace": namespace,
			"service":   serviceName,
		}).
		Get(client.APIurl + servicePath)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return model.Service{}, err
	}
	return *resp.Result().(*model.Service), nil
}

// GetServiceList -- consumes a namespace name
// returns a slice of Services OR a nil slice AND an error
func (client *Client) GetServiceList(namespace string) ([]model.Service, error) {
	resp, err := client.Request.
		SetResult([]model.Service{}).
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).
		Get(client.APIurl + servicesPath)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.Service), nil
}

// CreateService -- consumes a namespace name and a Service data,
// returns the created Service AND nil OR an uninitialized Service AND an error
func (client *Client) CreateService(namespace string, service model.Service) (model.Service, error) {
	resp, err := client.Request.
		SetResult(model.Service{}).
		SetBody(service).
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).Post(client.ResourceAddr + servicesPath)
	if err != nil {
		return model.Service{}, err
	}
	return *resp.Result().(*model.Service), nil
}

// DeleteService -- consumes a namespace, a servicce name,
// returns error in case of problem
func (client *Client) DeleteService(namespace, serviceName string) error {
	// TODO: check if errors with code > 399 are stored in err
	_, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
			"service":   serviceName,
		}).Delete(client.ResourceAddr + servicePath)
	return err
}

// UpdateService -- consumes a namespace, a service data,
// returns an ipdated Service OR an uninitialized Service AND an error
func (client *Client) UpdateService(namespace string, service model.Service) (model.Service, error) {
	resp, err := client.Request.
		SetResult(model.Service{}).
		SetBody(service).
		SetPathParams(map[string]string{
			"namespace": namespace,
			"service":   service.Name,
		}).Put(client.ResourceAddr + servicePath)
	if err != nil {
		return model.Service{}, err
	}
	return *resp.Result().(*model.Service), nil
}
