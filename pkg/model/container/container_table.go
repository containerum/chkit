package container

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model"
)

func (cont Container) RenderTable() string {
	return model.RenderTable(cont)
}

func (Container) TableHeaders() []string {
	return []string{"Name", "Image", "Limits"}
}

func (cont Container) TableRows() [][]string {
	return [][]string{{
		cont.Name,
		cont.Image,
		fmt.Sprintf("CPU: %4d mCPU\nMEM: %4d Mb", cont.Limits.CPU, cont.Limits.Memory),
	}}
}
