package namespace

import (
	"bytes"
	"time"

	"github.com/sirupsen/logrus"

	kubeModel "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/volume"
	"github.com/olekukonko/tablewriter"
)

type Namespace struct {
	CreatedAt *time.Time
	Label     string
	Access    string
	Volumes   []volume.Volume
	Resources kubeModel.Resources
	origin    kubeModel.Namespace
}

func NamespaceFromKube(kubeNameSpace kubeModel.Namespace) Namespace {
	ns := Namespace{
		Label:     kubeNameSpace.Label,
		Access:    kubeNameSpace.Access,
		Resources: kubeNameSpace.Resources,
		origin:    kubeNameSpace,
	}
	if kubeNameSpace.CreatedAt != nil {
		createdAt, err := time.Parse(model.TimestampFormat, *kubeNameSpace.CreatedAt)
		if err != nil {
			logrus.WithError(err).Debugf("invalid created_at timestamp")
		} else {
			ns.CreatedAt = &createdAt
		}
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
