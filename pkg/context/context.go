package context

import (
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/coblog"
)

type Context struct {
	Log                coblog.Log
	Version            string
	ConfigPath         string
	ConfigDir          string
	Namespace          Namespace
	Quiet              bool
	Changed            bool
	Client             chClient.Client
	AllowSelfSignedTLS bool
}

func (ctx *Context) StartCommand(command string) {
	ctx.Log.FieldLogger = ctx.Log.FieldLogger.WithField("command", command)
}

func (ctx *Context) ExitCommand() {
	ctx.Log.FieldLogger = ctx.Log.FieldLogger.WithField("command", nil)
}

func (ctx *Context) GetClient() *chClient.Client {
	return &ctx.Client
}

func (ctx *Context) SetAPI(api string) *Context {
	ctx.Client.APIaddr = api
	ctx.Changed = true
	return ctx
}

func (ctx *Context) SetNamespace(ns namespace.Namespace) *Context {
	ctx.Namespace = NamespaceFromModel(ns)
	ctx.Changed = true
	return ctx
}

type Storable struct {
	Namespace          Namespace
	Username           string
	Password           string
	API                string
	Version            string
	AllowSelfSignedTLS bool
}

func (config Storable) Merge(upd Storable) Storable {
	if upd.Namespace.Label != "" {
		config.Namespace.Label = upd.Namespace.Label
	}
	if upd.Namespace.ID != "" {
		config.Namespace.ID = upd.Namespace.ID
	}
	if upd.Namespace.OwnerLogin != "" {
		config.Namespace.OwnerLogin = upd.Namespace.OwnerLogin
	}
	if upd.API != "" {
		config.API = upd.API
	}
	if upd.Password != "" {
		config.Password = upd.Password
	}
	if upd.Username != "" {
		config.Username = upd.Username
	}
	config.AllowSelfSignedTLS = upd.AllowSelfSignedTLS
	return config
}

func (ctx *Context) GetStorable() Storable {
	return Storable{
		Version:            ctx.Version,
		Namespace:          ctx.Namespace,
		Username:           ctx.Client.Username,
		Password:           ctx.Client.Password,
		API:                ctx.Client.APIaddr,
		AllowSelfSignedTLS: ctx.AllowSelfSignedTLS,
	}
}

func (ctx *Context) SetStorable(config Storable) (configVersion string) {
	ctx.Namespace = config.Namespace
	ctx.Client.UserInfo = model.UserInfo{
		Username: config.Username,
		Password: config.Password,
	}
	if config.API != "" {
		ctx.Client.APIaddr = config.API
	}
	ctx.AllowSelfSignedTLS = config.AllowSelfSignedTLS
	return config.Version
}
