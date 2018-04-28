package volume

import (
	"time"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Volume struct {
	Label     string
	CreatedAt time.Time
	Access    string
	Replicas  uint
	Storage   uint
	origin    kubeModels.Volume
}

func VolumeFromKube(kv kubeModels.Volume) Volume {
	volume := Volume{
		Label:     kv.Label,
		CreatedAt: kv.CreateTime,
		Access:    kv.Access,
		Replicas:  uint(kv.Replicas),
		Storage:   uint(kv.Storage),
		origin:    kv,
	}
	return volume
}
