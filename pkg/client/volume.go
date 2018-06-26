package chClient

import (
	"github.com/containerum/chkit/pkg/model/volume"
	"github.com/containerum/chkit/pkg/util/coblog"
)

func (client *Client) GetVolume(namespaceID, volumeName string) (volume.Volume, error) {
	var vol volume.Volume
	var logger = coblog.Std.Component("GetVolume")
	err := retry(4, func() (bool, error) {
		kubeVolume, err := client.kubeAPIClient.GetVolume(namespaceID, volumeName)
		if err == nil {
			vol = volume.VolumeFromKube(kubeVolume)
		}
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logger.WithError(err).WithField("namespace", namespaceID).
			Errorf("unable to get volume %q", volumeName)
	}
	return vol, err
}

func (client *Client) GetVolumeList(namespaceID string) (volume.VolumeList, error) {
	var list volume.VolumeList
	var logger = coblog.Std.Component("GetVolumeList")
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetVolumeList(namespaceID)
		if err == nil {
			list = volume.VolumeListFromKube(kubeList)
		}
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logger.WithError(err).WithField("namespace", namespaceID).
			Errorf("unable to get volume list")
	}
	return list, err
}

func (client *Client) DeleteVolume(namespaceID, volumeName string) error {
	var logger = coblog.Std.Component("DeleteVolume")
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteVolume(namespaceID, volumeName)
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logger.WithError(err).WithField("namespace", namespaceID).
			Errorf("unable to delete volume")
	}
	return err
}
