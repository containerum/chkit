package volume

import (
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
)

type Volume struct {
	Label     string
	CreatedAt time.Time
	Access    string
	Replicas  uint
	Storage   uint
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
