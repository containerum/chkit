package volume

import (
	"time"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Volume struct {
	ID        string
	Label     string
	CreatedAt time.Time
	Access    string
	Replicas  uint
	Capacity  uint
	origin    kubeModels.Volume
}

func VolumeFromKube(kv kubeModels.Volume) Volume {
	volume := Volume{
		ID:        kv.ID,
		Label:     kv.Label,
		CreatedAt: kv.CreateTime,
		Access:    kv.Access,
		Replicas:  uint(kv.Replicas),
		Capacity:  uint(kv.Capacity),
		origin:    kv,
	}
	return volume
}
