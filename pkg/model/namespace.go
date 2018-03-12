package model

import (
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
)

type Namespace struct {
	CreatedAt *time.Time
	Label     string
	Access    string
	Volumes   []Volume
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
	ns.Volumes = make([]Volume, 0, len(kubeNameSpace.Volumes))
	for _, volume := range kubeNameSpace.Volumes {
		ns.Volumes = append(ns.Volumes,
			VolumeFromKube(volume))
	}
	return ns
}
