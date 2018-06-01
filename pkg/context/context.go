package context

import (
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
)

type Context struct {
	Version            string
	ConfigPath         string
	ConfigDir          string
	Namespace          Namespace
	Quiet              bool
	Changed            bool
	Client             chClient.Client
	AllowSelfSignedTLS bool
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
	AllowSelfSignedTLS bool
}

func (config Storable) Merge(upd Storable) Storable {
	if !upd.Namespace.IsEmpty() {
		config.Namespace = upd.Namespace
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
		Namespace:          ctx.Namespace,
		Username:           ctx.Client.Username,
		Password:           ctx.Client.Password,
		API:                ctx.Client.APIaddr,
		AllowSelfSignedTLS: ctx.AllowSelfSignedTLS,
	}
}

func (ctx *Context) SetStorable(config Storable) {
	ctx.Namespace = config.Namespace
	ctx.Client.UserInfo = model.UserInfo{
		Username: config.Username,
		Password: config.Password,
	}
	if config.API != "" {
		ctx.Client.APIaddr = config.API
	}
	ctx.AllowSelfSignedTLS = config.AllowSelfSignedTLS
}
