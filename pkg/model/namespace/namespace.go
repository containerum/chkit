package namespace

import (
	"bytes"
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/volume"
	"github.com/olekukonko/tablewriter"
)

var (
	_ model.TableRenderer = &Namespace{}
	_ model.TableRenderer = &NamespaceList{}
)

type NamespaceList []Namespace

func NamespaceListFromKube(kubeList []kubeModels.Namespace) NamespaceList {
	var list NamespaceList = make([]Namespace, 0, len(kubeList))
	for _, namespace := range kubeList {
		list = append(list, NamespaceFromKube(namespace))
	}
	return list
}

func (_ NamespaceList) TableHeaders() []string {
	return new(Namespace).TableHeaders()
}

func (list NamespaceList) TableRows() [][]string {
	row := make([][]string, 0, len(list))
	for _, ns := range list {
		row = append(row, ns.TableRows()...)
	}
	return row
}

type Namespace struct {
	CreatedAt *time.Time
	Label     string
	Access    string
	Volumes   []volume.Volume
}

func (_ *Namespace) TableHeaders() []string {
	return []string{"Label", "Created" /* "Access",*/, "Volumes"}
}

func (namespace *Namespace) TableRows() [][]string {
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

func (namespace *Namespace) RenderTable() string {
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

func (ns *Namespace) RenderVolumes() string {
	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader(new(volume.Volume).TableHeaders())
	for _, volume := range ns.Volumes {
		table.AppendBulk(volume.TableRows())
	}
	table.Render()
	return buf.String()
}
