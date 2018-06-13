package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/volume"
	"github.com/containerum/chkit/pkg/util/coblog"
)

func (client *Client) GetVolume(namespaceID, volumeName string) (volume.Volume, error) {
	var vol volume.Volume
	var logger = coblog.Std.Component("GetVolume")
	err := retry(4, func() (bool, error) {
		kubeVolume, err := client.kubeAPIClient.GetVolume(namespaceID, volumeName)
		switch {
		case err == nil:
			vol = volume.VolumeFromKube(kubeVolume)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource(),
			kubeErrors.ErrInternalError()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
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
		switch {
		case err == nil:
			list = volume.VolumeListFromKube(kubeList)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
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
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource(),
			kubeErrors.ErrInternalError()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logger.WithError(err).WithField("namespace", namespaceID).
			Errorf("unable to delete volume")
	}
	return err
}
