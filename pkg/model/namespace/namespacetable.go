package namespace

import (
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/volume"
)

var (
	_ model.TableRenderer = &Namespace{}
)

func (_ Namespace) TableHeaders() []string {
	return []string{"Label", "Created" /* "Access",*/, "Volumes"}
}

func (namespace Namespace) TableRows() [][]string {
	creationTime := ""
	if namespace.CreatedAt != nil {
		creationTime = namespace.CreatedAt.Format(model.CreationTimeFormat)
	}
	volumes := ""
	for i, volume := range namespace.Volumes {
		if i > 0 {
			volumes += "\n" + volume.Label
		}
		volumes += volume.Label
	}
	return [][]string{{
		namespace.Label,
		creationTime,
		//namespace.Access,
		volumes,
	}}
}

func (namespace Namespace) RenderTable() string {
	return model.RenderTable(namespace)
}
func NamespaceFromKube(kubeNameSpace kubeModels.Namespace) Namespace {
	ns := Namespace{
		Label:  kubeNameSpace.Label,
		Access: kubeNameSpace.Access,
	}
	if kubeNameSpace.CreatedAt != nil {
		createdAt := time.Unix(*kubeNameSpace.CreatedAt, 0)
		ns.CreatedAt = &createdAt
	}
	ns.Volumes = make([]volume.Volume, 0, len(kubeNameSpace.Volumes))
	for _, kubeVolume := range kubeNameSpace.Volumes {
		ns.Volumes = append(ns.Volumes,
			volume.VolumeFromKube(kubeVolume))
	}
	return ns
}
