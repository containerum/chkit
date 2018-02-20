package cmd

import (
	"fmt"
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	kubeAPIpodRootPath = "/namespaces/{namespace}/pods"
	kubeAPIpodPath     = "/namespaces/{namespace}/pods/{pod}"
)

// DeletePod -- deletes pod in provided namespace
func (client *Client) DeletePod(namespace, pod string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"pod": pod,
		}).
		Delete(client.APIurl + kubeAPIpodPath)
	if err != nil {
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return nil
	default:
		return fmt.Errorf("%s", string(resp.Body()))
	}
}

// GetPod -- gets a particular pod by name.
func (client *Client) GetPod(namespace, pod string) (model.Pod, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
			"pod":       pod,
		}).
		Get(client.APIurl + kubeAPIpodPath)
	if err != nil {
		return model.Pod{}, err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return *resp.Result().(*model.Pod), nil
	default:
		return model.Pod{}, fmt.Errorf("%s", string(resp.Body()))
	}
}

func (client *Client) GetPodList(namespace string) ([]model.Pod, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).
		Get(client.APIurl + kubeAPIpodRootPath)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return *resp.Result().(*[]model.Pod), nil
	default:
		return nil, fmt.Errorf("%s", string(resp.Body()))
	}
}
