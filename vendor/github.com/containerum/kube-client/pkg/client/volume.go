package client

import (
	"github.com/containerum/kube-client/pkg/rest"

	"github.com/containerum/kube-client/pkg/model"
)

const (
	volumesPath = "/namespaces/{namespace}/volumes"
	volumePath  = "/namespaces/{namespace}/volumes/{volume}"
)

// DeleteVolume -- deletes Volume with provided volume name
func (client *Client) DeleteVolume(namespace, volumeName string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: volumePath,
			Params: rest.P{
				"namespace": namespace,
				"volume":    volumeName,
			},
		},
	})
}

// GetVolume -- return User Volume by name,
// consumes optional userID param
func (client *Client) GetVolume(namespace, volumeName string) (model.Volume, error) {
	var volume model.Volume
	err := client.RestAPI.Get(rest.Rq{
		Result: &volume,
		URL: rest.URL{
			Path: volumePath,
			Params: rest.P{
				"namespace": namespace,
				"volume":    volumeName,
			},
		},
	})
	return volume, err
}

// GetVolumeList -- get list of volumes,
// consumes optional user ID and filter parameters.
// Returns new_access_level as access if user role = user.
func (client *Client) GetVolumeList(namespace string) (model.VolumesList, error) {
	var volumeList model.VolumesList
	err := client.RestAPI.Get(rest.Rq{
		Result: &volumeList,
		URL: rest.URL{
			Path: volumesPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return volumeList, err
}
