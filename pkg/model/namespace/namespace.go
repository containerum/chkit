package namespace

import (
	"bytes"
	"time"

	"github.com/containerum/chkit/pkg/model/volume"
	"github.com/olekukonko/tablewriter"
)

type Namespace struct {
	CreatedAt *time.Time
	Label     string
	Access    string
	Volumes   []volume.Volume
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
