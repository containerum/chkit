package model

import (
	"fmt"
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
)

const (
	CreationTimeFormat = "2 Jan 2006 15:04 -0700 MST "
)

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
		volume.CreatedAt.Format(CreationTimeFormat),
		volume.Access,
		fmt.Sprintf("%d", volume.Replicas),
		fmt.Sprintf("%d", volume.Storage),
	}}
}

func (volume *Volume) RenderTable() string {
	return RenderTable(volume)
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
