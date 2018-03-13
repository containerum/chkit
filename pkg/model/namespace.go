package model

import (
	"bytes"
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/olekukonko/tablewriter"
)

type Namespace struct {
	CreatedAt *time.Time
	Label     string
	Access    string
	Volumes   []Volume
}

func (_ *Namespace) TableHeaders() []string {
	return []string{"Label", "Created", "Access", "Volumes"}
}

func (namespace *Namespace) TableRows() [][]string {
	creationTime := "unknown"
	if namespace.CreatedAt != nil {
		creationTime = namespace.CreatedAt.Format(CreationTimeFormat)
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
		namespace.Access,
		volumes,
	}}
}

func (namespace *Namespace) RenderTable() string {
	return RenderTable(namespace)
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

func (ns *Namespace) RenderVolumes() string {
	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader(new(Volume).TableHeaders())
	for _, volume := range ns.Volumes {
		table.AppendBulk(volume.TableRows())
	}
	table.Render()
	return buf.String()
}
