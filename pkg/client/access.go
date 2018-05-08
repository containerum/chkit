package chClient

import "github.com/containerum/chkit/pkg/model/access"

func (client *Client) GetAccess(nsName string) (access.Access, error) {
	ns, err := client.GetNamespace(nsName)
	return access.AccessFromNamespace(ns), err
}

func (client *Client) GetAccessList() (access.AccessList, error) {
	list, err := client.GetNamespaceList()
	return access.AccessListFromNamespaces(list), err
}
