package chClient

import (
	"io"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/client"
	"github.com/containerum/chkit/pkg/model/pod"
)

func (client *Client) GetPod(ns, podname string) (pod.Pod, error) {
	var gainedPod pod.Pod
	err := retry(4, func() (bool, error) {
		kubePod, err := client.kubeAPIClient.GetPod(ns, podname)
		switch {
		case err == nil:
			gainedPod = pod.PodFromKube(kubePod)
			return false, nil
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return gainedPod, err
}

func (client *Client) GetPodList(ns string) (pod.PodList, error) {
	var gainedList pod.PodList
	err := retry(4, func() (bool, error) {
		kubePod, err := client.kubeAPIClient.GetPodList(ns)
		switch {
		case err == nil:
			gainedList = pod.PodListFromKube(kubePod)
			return false, nil
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return gainedList, err
}

type GetPodLogsParams = client.GetPodLogsParams

func (client *Client) GetPodLogs(params client.GetPodLogsParams) (io.ReadCloser, error) {
	var rc io.ReadCloser
	err := retry(4, func() (bool, error) {
		var getErr error
		rc, getErr = client.kubeAPIClient.GetPodLogs(params)
		return false, getErr
	})
	return rc, err
}
