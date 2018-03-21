package pod

import (
	"strconv"
	"strings"
	"time"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = new(Pod)
)

func (pod *Pod) RenderTable() string {
	return model.RenderTable(pod)
}

func (_ *Pod) TableHeaders() []string {
	return []string{"Name", "Host", "Phase", "Restarts", "StartedAt", "Containers"}
}

func (pod *Pod) TableRows() [][]string {
	return [][]string{
		{
			pod.Name,
			pod.Hostname,
			pod.Status.Phase,
			strconv.Itoa(pod.Status.RestartCount),
			pod.Status.StartedAt.Format(time.RFC1123),
			strings.Join(pod.Containers, "\n"),
		},
	}
}
