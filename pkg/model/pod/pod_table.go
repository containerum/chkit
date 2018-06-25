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

func (Pod) TableHeaders() []string {
	return []string{"Label", "Status", "Restarts", "Age"}
}

func (pod Pod) TableRows() [][]string {
	age := "unknown"
	if pod.CreatedAt.Unix() != 0 {
		age = model.Age(pod.CreatedAt)
	}
	return [][]string{{
		pod.Name,
		pod.Status.Phase,
		strconv.Itoa(pod.Status.RestartCount),
		age,
	}}
}
