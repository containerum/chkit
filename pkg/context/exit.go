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
	var logger = ctx.Log.Component(fmt.Sprintf("%v postrun", command.CommandPath()))
	logger.Debugf("START")
	defer logger.Debugf("END")
	ctx.RunDeffered()
}
