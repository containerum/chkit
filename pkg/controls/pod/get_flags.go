package pod

import (
	"strings"

	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/pod"
	"github.com/ninedraft/boxofstuff/str"
)

type GetFlags struct {
	Status   string `desc:"include in result pods with custom status"`
	Failed   bool   `desc:"include in result pods with status 'Failed'"`
	Pending  bool   `desc:"include in result pods with status 'Pending'"`
	Running  bool   `desc:"include in result pods with status 'Running'"`
	Filename string `desc:"output filename, STDOUT if empty of '-'"`
	Output   string `desc:"output format, yaml/json" flag:"output o"`
}

func (gf GetFlags) StatusFilter() func(po pod.Pod) bool {
	var statuses = gf.Statuses()
	return func(po pod.Pod) bool {
		return statuses.Contains(strings.ToLower(po.Status.Phase))
	}
}

func (gf GetFlags) IsStatusesDefined() bool {
	return gf.Running || gf.Pending || gf.Failed || gf.Status != ""
}

func (gf GetFlags) ExportConfig() export.ExportConfig {
	return export.ExportConfig{
		Filename: gf.Filename,
		Format:   export.ExportFormat(gf.Output),
	}
}

func (gf GetFlags) Statuses() str.Vector {
	if gf.IsStatusesDefined() {
		var statuses = str.Vector{gf.Status}
		if gf.Running {
			statuses = append(statuses, "running")
		}
		if gf.Pending {
			statuses = append(statuses, "pending")
		}
		if gf.Failed {
			statuses = append(statuses, "failed")
		}
		return statuses.Filter(str.Longer(0))
	}
	return []string{"running", "pending", "failed"}
}
