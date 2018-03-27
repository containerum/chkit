package pod

import (
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = new(Pod)
)

func (pod *Pod) RenderTable() string {
	return model.RenderTable(pod)
}

func (_ *Pod) TableHeaders() []string {
	return []string{"Name", "Host", "Phase", "Restarts", "Age", "Containers"}
}

func (pod *Pod) TableRows() [][]string {
	return [][]string{
		{
			pod.Name,
			pod.Hostname,
			pod.Status.Phase,
			strconv.Itoa(pod.Status.RestartCount),
			model.TimestampFormat(pod.Status.StartedAt),
			strings.Join(pod.Containers, "\n"),
		},
	}
}
