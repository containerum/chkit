package context

import (
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
)

type Context struct {
	Version    string
	ConfigPath string
	ConfigDir  string
	Namespace  string
	Quiet      bool
	Changed    bool
	Client     chClient.Client
}

type Storable struct {
	Namespace string
	Username  string
	Password  string
}

func (ctx *Context) GetStorable() Storable {
	return Storable{
		Namespace: ctx.Namespace,
		Username:  ctx.Client.Username,
		Password:  ctx.Client.Password,
	}
}

func (ctx *Context) SetStorable(config Storable) {
	ctx.Namespace = config.Namespace
	ctx.Client.UserInfo = model.UserInfo{
		Username: config.Username,
		Password: config.Password,
	}
}
