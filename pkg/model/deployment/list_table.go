package deployment

import "github.com/containerum/chkit/pkg/model"

var (
	_ model.TableRenderer = DeploymentList{}
)

func (list DeploymentList) RenderTable() string {
	return model.RenderTable(list)
}

func (_ DeploymentList) TableHeaders() []string {
	return new(Deployment).TableHeaders()
}

func (list DeploymentList) TableRows() [][]string {
	table := make([][]string, 0, len(list))
	for _, depl := range list {
		table = append(table, depl.TableRows()...)
	}
	return table
}
