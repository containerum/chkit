package chClient

import (
	"io"

	"github.com/containerum/chkit/pkg/model/pod"
	"github.com/containerum/kube-client/pkg/client"
)

func (client *Client) GetPod(ns, podname string) (pod.Pod, error) {
	var gainedPod pod.Pod
	err := retry(4, func() (bool, error) {
		kubePod, err := client.kubeAPIClient.GetPod(ns, podname)
		if err == nil {
			gainedPod = pod.PodFromKube(kubePod)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedPod, err
}

func (client *Client) GetPodList(ns string) (pod.PodList, error) {
	var gainedList pod.PodList
	err := retry(4, func() (bool, error) {
		kubePod, err := client.kubeAPIClient.GetPodList(ns)
		if err == nil {
			gainedList = pod.PodListFromKube(kubePod)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedList, err
}

func (client *Client) GetDeploymentPodList(ns, deploy string) (pod.PodList, error) {
	var gainedList pod.PodList
	err := retry(4, func() (bool, error) {
		kubePod, err := client.kubeAPIClient.GetDeploymentPodList(ns, deploy)
		if err == nil {
			gainedList = pod.PodListFromKube(kubePod)
		}
		return HandleErrorRetry(client, err)
	})
	return gainedList, err
}

func (client *Client) DeletePod(namespace, pod string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeletePod(namespace, pod)
		return HandleErrorRetry(client, err)
	})
}

type GetPodLogsParams = client.GetPodLogsParams

func (client *Client) GetPodLogs(params client.GetPodLogsParams) (io.ReadCloser, error) {
	var rc io.ReadCloser
	err := retry(4, func() (bool, error) {
		var getErr error
		rc, getErr = client.kubeAPIClient.GetPodLogs(params)
		return HandleErrorRetry(client, getErr)
	})
	return rc, err
}
