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
	namespace          Namespace
	Changed            bool
	client             chClient.Client
	allowSelfSignedTLS bool

	deferred []func()
}

func (ctx *Context) GetNamespace() Namespace {
	return ctx.namespace
}

func (ctx *Context) SetNamespace(ns namespace.Namespace) *Context {
	ctx.namespace = NamespaceFromModel(ns)
	ctx.Changed = true
	return ctx
}

func (ctx *Context) SetTemporaryNamespace(ns namespace.Namespace) *Context {
	ctx.namespace = NamespaceFromModel(ns)
	return ctx
}

func (ctx *Context) GetSelfSignedTLSRule() bool {
	return ctx.allowSelfSignedTLS
}

func (ctx *Context) SetSelfSignedTLSRule(allow bool) *Context {
	ctx.allowSelfSignedTLS = allow
	ctx.Changed = true
	return ctx
}

func (ctx *Context) GetClient() *chClient.Client {
	return &ctx.client
}

func (ctx *Context) SetAPI(api string) *Context {
	ctx.client.APIaddr = api
	ctx.Changed = true
	return ctx
}

func (ctx *Context) GetAPI() string {
	return ctx.client.APIaddr
}

func (ctx *Context) SetAuth(login, password string) *Context {
	ctx.client.Username = login
	ctx.client.Password = password
	ctx.Changed = true
	return ctx
}

func (ctx *Context) GetAuth() model.UserInfo {
	return ctx.client.UserInfo
}

func (ctx *Context) StartCommand(command string) {
	ctx.Log.FieldLogger = ctx.Log.FieldLogger.WithField("command", command)
}

func (ctx *Context) ExitCommand() {
	ctx.Log.FieldLogger = ctx.Log.FieldLogger.WithField("command", nil)
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
		Namespace:          ctx.namespace,
		Username:           ctx.client.Username,
		Password:           ctx.client.Password,
		API:                ctx.client.APIaddr,
		AllowSelfSignedTLS: ctx.allowSelfSignedTLS,
	}
}

func (ctx *Context) SetStorable(config Storable) (configVersion string) {
	ctx.namespace = config.Namespace
	ctx.client.UserInfo = model.UserInfo{
		Username: config.Username,
		Password: config.Password,
	}
	if config.API != "" {
		ctx.client.APIaddr = config.API
	}
	ctx.allowSelfSignedTLS = config.AllowSelfSignedTLS
	return config.Version
}
