package context

import (
	"os"

	"github.com/spf13/cobra"
)

func (ctx *Context) Defer(f func()) *Context {
	ctx.deferred = append(ctx.deferred, f)
	return ctx
}

func (ctx *Context) Exit(code int) {
	ctx.RunDeffered()
	os.Exit(code)
}

func (ctx *Context) RunDeffered() *Context {
	for _, deferred := range ctx.deferred {
		deferred()
	}
	return ctx
}

func (ctx *Context) CobraPostrun(command *cobra.Command, args []string) {
	ctx.RunDeffered()
}
