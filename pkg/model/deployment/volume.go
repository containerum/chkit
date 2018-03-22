package deployment

import (
	"fmt"
	"strings"
	"time"

	"github.com/containerum/chkit/pkg/model"
)

type Volume struct {
	Label     string
	CreatedAt time.Time
	Access    string
	Storage   int
	Replicas  int
}

func (volume *Volume) String() string {
	return strings.Join([]string{
		"Label: " + volume.Label,
		"Created at" + volume.CreatedAt.Format(model.CreationTimeFormat),
		"Access: " + volume.Access,
		"Storage: " + fmt.Sprintf("%dGb", volume.Storage),
		"Replicas: " + fmt.Sprintf("%d", volume.Replicas),
	}, "\n")
}
