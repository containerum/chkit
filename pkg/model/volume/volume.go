package volume

import (
	"fmt"
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &Volume{}
	_ model.TableRenderer = &VolumeList{}
)

type VolumeList []Volume

func VolumeListFromKube(kubeList []kubeModels.Volume) VolumeList {
	var list VolumeList = make([]Volume, 0, len(kubeList))
	for _, volume := range kubeList {
		list = append(list, VolumeFromKube(volume))
	}
	return list
}

func (_ VolumeList) TableHeaders() []string {
	return new(Volume).TableHeaders()
}

func (list VolumeList) TableRows() [][]string {
	rows := make([][]string, 0, len(list))
	for _, volume := range list {
		rows = append(rows, volume.TableRows()...)
	}
	return rows
}

type Volume struct {
	Label     string
	CreatedAt time.Time
	Access    string
	Replicas  uint
	Storage   uint
}

func (_ *Volume) TableHeaders() []string {
	return []string{"Label", "Created", "Access", "Replicas", "Storage, GB"}
}

func (volume *Volume) TableRows() [][]string {
	return [][]string{{
		volume.Label,
		volume.CreatedAt.Format(model.CreationTimeFormat),
		volume.Access,
		fmt.Sprintf("%d", volume.Replicas),
		fmt.Sprintf("%d", volume.Storage),
	}}
}

func (volume *Volume) RenderTable() string {
	return model.RenderTable(volume)
}
func VolumeFromKube(kv kubeModels.Volume) Volume {
	volume := Volume{
		Label:     kv.Label,
		CreatedAt: kv.CreateTime,
		Access:    kv.Access,
		Replicas:  uint(kv.Replicas),
		Storage:   uint(kv.Storage),
	}
	return volume
}
