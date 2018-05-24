package namespace

import (
	"testing"
	"time"

	"github.com/containerum/chkit/pkg/model/volume"
)

func TestNamespaceRenderToTable(test *testing.T) {
	creationTime := time.Now()
	ns := Namespace{
		Label:     "mushrooms",
		Access:    "r-only",
		CreatedAt: &creationTime,
		Volumes: []volume.Volume{
			{
				ID:        "id1",
				Label:     "newton",
				CreatedAt: time.Now(),
				Access:    "r/w",
				Replicas:  10,
				Capacity:  5,
			},
			{
				ID:        "id2",
				Label:     "max",
				CreatedAt: time.Now(),
				Access:    "r",
				Replicas:  4,
				Capacity:  10,
			},
		},
	}
	test.Logf("\n%v", ns.RenderTable())
	test.Logf("\n%v", ns.RenderVolumes())
}
