package volume

import (
	"time"

	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

var (
	_ model.Renderer = Volume{}
)

type Volume kubeModels.Volume

func VolumeFromKube(kv kubeModels.Volume) Volume {
	return Volume(kv).Copy()
}

func (volume Volume) ToKube() kubeModels.Volume {
	return kubeModels.Volume(volume.Copy())
}

func (volume Volume) Age() string {
	if volume.CreatedAt == nil {
		return "undefined"
	}
	var timestamp, _ = time.Parse(*volume.CreatedAt, model.TimestampFormat)
	return model.Age(timestamp)
}

func (volume Volume) Copy() Volume {
	var cp = volume
	cp.Users = append(make([]kubeModels.UserAccess, 0, len(volume.Users)), volume.Users...)
	return cp
}

func (volume Volume) UserNames() []string {
	var names = make([]string, 0, len(volume.Users))
	for _, user := range volume.Users {
		names = append(names, user.Username)
	}
	return names
}

func (volume Volume) String() string {
	return volume.OwnerAndName()
}

func (volume Volume) OwnerAndName() string {
	return volume.OwnerLogin + "/" + volume.Name
}
