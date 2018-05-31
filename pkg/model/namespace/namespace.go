package namespace

import (
	"time"

	"fmt"

	"github.com/containerum/chkit/pkg/model"
	kubeModel "github.com/containerum/kube-client/pkg/model"
)

type Namespace kubeModel.Namespace

func NamespaceFromKube(kubeNameSpace kubeModel.Namespace) Namespace {
	return Namespace(kubeNameSpace).Copy()
}

func (namespace Namespace) Copy() Namespace {
	namespace.Users = append(make([]kubeModel.UserAccess, 0, len(namespace.Users)), namespace.Users...)
	return namespace
}

func (namespace Namespace) ToKube() kubeModel.Namespace {
	return kubeModel.Namespace(namespace.Copy())
}

func (namespace Namespace) UserNames() []string {
	var names = make([]string, 0, len(namespace.Users))
	for _, user := range namespace.Users {
		names = append(names, user.Username)
	}
	return names
}

func (namespace Namespace) Age() string {
	if namespace.CreatedAt == nil {
		return "undefined"
	}
	var timestamp, _ = time.Parse(model.TimestampFormat, *namespace.CreatedAt)
	return model.Age(timestamp)
}

func (namespace Namespace) UsageCPU() string {
	var hard, used = namespace.Resources.Hard, namespace.Resources.Used
	if used == nil {
		return fmt.Sprintf("%d mCPU", hard.CPU)
	}
	return fmt.Sprintf("%d/%d mCPU", used.CPU, hard.CPU)
}

func (namespace Namespace) UsageMemory() string {
	var hard, used = namespace.Resources.Hard, namespace.Resources.Used
	if used == nil {
		return fmt.Sprintf("%d Mb", hard.Memory)
	}
	return fmt.Sprintf("%d/%d Mb", used.Memory, hard.Memory)
}

func (namespace Namespace) LabelAndID() string {
	return fmt.Sprintf("%s (%s)", namespace.Label, namespace.OwnerLogin)
}
