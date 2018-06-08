package login

import (
	"os"

	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

type Flags struct {
	Username  string `flag:"username u"`
	Password  string `flag:"password p"`
	Namespace string `flag:"-"`
}

func Login(ctx *context.Context) *cobra.Command {
	var flags Flags
	command := &cobra.Command{
		Use:   "login",
		Short: "Login to system",
		Run: func(command *cobra.Command, args []string) {
			if err := clisetup.SetupLogs(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			flags.Namespace, _ = command.Flags().GetString("namespace")
			if err := RunLogin(ctx, flags); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			ctx.Log.Command("login").Debugf("saving config")
			postrun.PostRun(ctx)
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}

func RunLogin(ctx *context.Context, flags Flags) error {
	var logger = ctx.Log.Component("RunLogin")
	logger.Debugf("start")
	defer logger.Debugf("end")
	ctx.Client.Username = flags.Username
	ctx.Client.Password = flags.Password
	ctx.Changed = true
	logger.Debugf("start app setup")
	if err := Setup(ctx); err != nil {
		angel.Angel(ctx, err)
		os.Exit(1)
	}
	logger.Debugf("end setup")

	switch flags.Namespace {
	case "-":
		clisetup.GetDefaultNS(ctx, true)
	case "":
		clisetup.GetDefaultNS(ctx, false)
	default:
		nsList, err := ctx.Client.GetNamespaceList()
		logger.Debugf("Getting namespace list")
		if err != nil {
			logger.WithError(err).Errorf("unable to get namespace lsit")
			fmt.Println(err)
			os.Exit(1)
		}
		ns, ok := nsList.GetByUserFriendlyID(flags.Namespace)
		if !ok {
			fmt.Printf("Namespace %q not found!\n%s",
				flags.Namespace, strings.Join(nsList.OwnersAndLabels(), "\n"))
			os.Exit(1)
		}
		ctx.SetNamespace(ns)
	}
	return nil
}
