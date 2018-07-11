package chClient

import (
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetIngress(ns, domain string) (ingress.Ingress, error) {
	var ingr ingress.Ingress
	err := retry(4, func() (bool, error) {
		kubeIngress, err := client.kubeAPIClient.GetIngress(ns, domain)
		if err == nil {
			ingr = ingress.IngressFromKube(kubeIngress)
		}
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to get ingress")
	}
	return ingr, err
}

func (client *Client) GetIngressList(ns string) (ingress.IngressList, error) {
	var list ingress.IngressList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetIngressList(ns)
		if err == nil {
			list = ingress.IngressListFromKube(kubeList)
		}
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to get ingress list")
	}
	return list, err
}

func (client *Client) CreateIngress(ns string, ingr ingress.Ingress) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.AddIngress(ns, ingr.ToKube())
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).
			WithField("namespace", ns).
			Errorf("unable to create ingress")
	}
	return err
}

func (client *Client) ReplaceIngress(ns string, ingr ingress.Ingress) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.UpdateIngress(ns, ingr.Name, ingr.ToKube())
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).
			WithField("namespace", ns).
			Errorf("unable to create ingress")
	}
	return err
}

func (client *Client) DeleteIngress(ns, name string) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteIngress(ns, name)
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to get ingress")
	}
	return err
}
