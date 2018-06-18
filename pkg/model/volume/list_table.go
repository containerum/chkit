package volume

import "github.com/containerum/chkit/pkg/model"

var (
	_ model.TableRenderer = &VolumeList{}
)

func (_ VolumeList) TableHeaders() []string {
	return new(Volume).TableHeaders()
}

func (list VolumeList) TableRows() [][]string {
	rows := make([][]string, 0, len(list))
	for _, volume := range list {
		rows = append(rows, volume.TableRows()...)
	}
	return rows
}

func (list VolumeList) RenderTable() string {
	return model.RenderTable(list)
}
