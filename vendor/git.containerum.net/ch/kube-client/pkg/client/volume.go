package client

import (
	"git.containerum.net/ch/kube-client/pkg/rest"

	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	resourceVolumeRootPath   = "/volume"
	resourceVolumePath       = "/volume/{volume}"
	resourceVolumeNamePath   = "/volume/{volume}/name"
	resourceVolumeAccessPath = "/volume/{volume}/access"
)

// DeleteVolume -- deletes Volume with provided volume name
func (client *Client) DeleteVolume(volumeName string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: client.APIurl + resourceVolumePath,
			Params: rest.P{
				"volume": volumeName,
			},
		},
	})
}

// GetVolume -- return User Volume by name,
// consumes optional userID param
func (client *Client) GetVolume(volumeName string) (model.Volume, error) {
	var volume model.Volume
	err := client.RestAPI.Get(rest.Rq{
		Result: &volume,
		URL: rest.URL{
			Path: client.APIurl + resourceVolumePath,
			Params: rest.P{
				"volume": volumeName,
			},
		},
	})
	return volume, err
}

// GetVolumeList -- get list of volumes,
// consumes optional user ID and filter parameters.
// Returns new_access_level as access if user role = user.
// Should have filters: not deleted, limited, not limited, owner, not owner.
func (client *Client) GetVolumeList(filter *string) ([]model.Volume, error) {
	var volumeList []model.Volume
	var query = rest.P{}
	if filter != nil {
		query["filter"] = *filter
	}
	err := client.RestAPI.Get(rest.Rq{
		Result: &volumeList,
		URL: rest.URL{
			Path:   client.APIurl + resourceVolumeRootPath,
			Params: rest.P{},
		},
	})
	return volumeList, err
}

//RenameVolume -- change volume name
func (client *Client) RenameVolume(volumeName, newName string) error {
	return client.RestAPI.Put(rest.Rq{
		Body: model.ResourceUpdateName{
			Label: newName,
		},
		URL: rest.URL{
			Path: client.APIurl + resourceVolumeNamePath,
			Params: rest.P{
				"volume": volumeName,
			},
		},
	})
}

// SetVolumeAccess -- sets User Volume access
func (client *Client) SetVolumeAccess(volumeName string, accessData model.ResourceUpdateUserAccess) error {
	return client.RestAPI.Post(rest.Rq{
		Body: accessData,
		URL: rest.URL{
			Path: client.APIurl + resourceVolumeAccessPath,
			Params: rest.P{
				"volume": volumeName,
			},
		},
	})
}

// DeleteAccess -- deletes user Volume access
func (client *Client) DeleteAccess(volumeName, username string) error {
	return client.RestAPI.Delete(rest.Rq{
		Body: model.ResourceUpdateUserAccess{
			Username: username,
		},
		URL: rest.URL{
			Path: client.APIurl + resourceVolumeAccessPath,
			Params: rest.P{
				"volume": volumeName,
			},
		},
	})
}
