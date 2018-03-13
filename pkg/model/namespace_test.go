package model

import (
	"testing"
	"time"
)

func TestNamespaceRenderToTable(test *testing.T) {
	ns := Namespace{
		Volumes: []Volume{
			{
				Label:     "newton",
				CreatedAt: time.Now(),
				Access:    "r/w",
				Replicas:  10,
				Storage:   5,
			},
		},
	}
	test.Logf("\n%v", ns.RenderTextTable())
}
