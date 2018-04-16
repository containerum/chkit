package context

import (
	"path"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/model"
)

var GlobalContext = context{
	Version:    semver.MustParse("3.0.1-alpha").String(),
	ConfigDir:  configdir.ConfigDir(),
	ConfigPath: path.Join(configdir.ConfigDir(), "config.toml"),
}

type context struct {
	Version     string
	ConfigPath  string
	ConfigDir   string
	Fingerprint string
	Namespace   string
	Quiet       bool
	Changed     bool
	Client      chClient.Client
}

type Storable struct {
	Namespace string
	Username  string
	Password  string
}

func (ctx *context) GetStorable() Storable {
	return Storable{
		Namespace: ctx.Namespace,
		Username:  ctx.Client.Username,
		Password:  ctx.Client.Password,
	}
}

func (ctx *context) SetStorable(config Storable) {
	ctx.Namespace = config.Namespace
	ctx.Client.UserInfo = model.UserInfo{
		Username: config.Username,
		Password: config.Password,
	}
}
