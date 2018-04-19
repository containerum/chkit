package pod

import (
	"strconv"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = new(Pod)
)

func (pod Pod) RenderTable() string {
	return model.RenderTable(pod)
}

func (_ Pod) TableHeaders() []string {
	return []string{"Label", "Host", "Status", "Restarts", "Age"}
}

func (pod Pod) TableRows() [][]string {
	age := "unknown"
	if pod.Status.StartedAt.Unix() != 0 {
		age = model.Age(pod.Status.StartedAt)
	}
	return [][]string{
		{
			pod.Name,
			pod.Hostname,
			pod.Status.Phase,
			strconv.Itoa(pod.Status.RestartCount),
			age,
		},
	}
}
