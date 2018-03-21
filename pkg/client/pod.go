package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"github.com/containerum/chkit/pkg/model/pod"
)

func (client *Client) GetPod(ns, podname string) (pod.Pod, error) {
	var gainedPod pod.Pod
	err := retry(4, func() error {
		kubePod, err := client.kubeAPIClient.GetPod(ns, podname)
		switch {
		case err == nil:
			gainedPod = pod.PodFromKube(kubePod)
			return nil
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return client.Auth()
		default:
			return err
		}
	})
	return gainedPod, err
}
