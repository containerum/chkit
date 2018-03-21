package namespace

import (
	"testing"
	"time"
)

func TestNamespaceRenderToTable(test *testing.T) {
	creationTime := time.Now()
	ns := Namespace{
		Label:     "mushrooms",
		Access:    "r-only",
		CreatedAt: &creationTime,
		Volumes: []Volume{
			{
				Label:     "newton",
				CreatedAt: time.Now(),
				Access:    "r/w",
				Replicas:  10,
				Storage:   5,
			},
			{
				Label:     "max",
				CreatedAt: time.Now(),
				Access:    "r",
				Replicas:  4,
				Storage:   10,
			},
		},
	}
	test.Logf("\n%v", ns.RenderTable())
	test.Logf("\n%v", ns.RenderVolumes())
}
