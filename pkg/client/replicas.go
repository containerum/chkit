package chClient

func (client *Client) SetReplicas(ns, depl string, n uint64) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.SetReplicas(ns, depl, int(n))
		return HandleErrorRetry(client, err)
	})
}
