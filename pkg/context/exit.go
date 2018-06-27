package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (ctx *Context) Defer(f func()) *Context {
	ctx.Log.Component("defer").Debugf("adding func")
	ctx.deferred = append(ctx.deferred, f)
	return ctx
}

func (ctx *Context) DeferCobra(f func(cmd *cobra.Command, args []string)) *Context {
	ctx.deferredCobra = append(ctx.deferredCobra, f)
	return ctx
}

func (ctx *Context) Exit(code int) {
	ctx.RunDeferred()
	os.Exit(code)
}

func (ctx *Context) RunDeferred() *Context {
	for _, deferred := range ctx.deferred {
		deferred()
	}
	return ctx
}

func (ctx *Context) CobraPostRun(cmd *cobra.Command, args []string) {
	var logger = ctx.Log.Component(fmt.Sprintf("%v postrun", cmd.CommandPath()))
	logger.Debugf("START")
	defer logger.Debugf("END")
	for _, f := range ctx.deferredCobra {
		f(cmd, args)
	}
	ctx.RunDeferred()
}
