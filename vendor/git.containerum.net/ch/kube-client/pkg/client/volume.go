package client

import (
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"
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
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"volume": volumeName,
		}).
		SetError(cherry.Err{}).
		Delete(client.ResourceAddr + resourceVolumePath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// GetVolume -- return User Volume by name,
// consumes optional userID param
func (client *Client) GetVolume(volumeName string) (model.Volume, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"volume": volumeName,
		}).
		SetError(cherry.Err{}).
		SetResult(model.Volume{}).
		Get(client.ResourceAddr + resourceVolumePath)
	if err = MapErrors(resp, err, http.StatusOK); err != nil {
		return model.Volume{}, err
	}
	return *resp.Result().(*model.Volume), nil
}

// GetVolumeList -- get list of volumes,
// consumes optional user ID and filter parameters.
// Returns new_access_level as access if user role = user.
// Should have filters: not deleted, limited, not limited, owner, not owner.
func (client *Client) GetVolumeList(filter *string) ([]model.Volume, error) {
	req := client.Request.
		SetResult([]model.Volume{}).
		SetError(cherry.Err{})
	if filter != nil {
		req.SetQueryParam("filter", *filter)
	}
	resp, err := req.Get(client.ResourceAddr + resourceVolumeRootPath)
	if err = MapErrors(resp, err, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.Volume), nil
}

//RenameVolume -- change volume name
func (client *Client) RenameVolume(volumeName, newName string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"volume": volumeName,
		}).
		SetError(cherry.Err{}).
		SetBody(model.ResourceUpdateName{Label: newName}).
		Put(client.ResourceAddr + resourceVolumeNamePath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// SetAccess -- sets User Volume access
func (client *Client) SetAccess(volumeName string, accessData model.ResourceUpdateUserAccess) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"volume": volumeName,
		}).
		SetError(cherry.Err{}).
		SetBody(accessData).
		Post(client.ResourceAddr + resourceVolumeAccessPath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// DeleteAccess -- deletes user Volume access
func (client *Client) DeleteAccess(volumeName, username string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"volume": volumeName,
		}).
		SetBody(model.ResourceUpdateUserAccess{
			Username: username,
		}).
		SetError(cherry.Err{}).
		Delete(client.ResourceAddr + resourceVolumeAccessPath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}
