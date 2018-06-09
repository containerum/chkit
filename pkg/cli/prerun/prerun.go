package prerun

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/login"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

func PreRun(ctx *context.Context) error {
	clisetup.SetupLogs(ctx)
	var logger = ctx.Log.Component("PreRun")
	err := configuration.SyncConfig(ctx)
	switch err {
	case nil:
		// pass
	case configuration.ErrIncompatibleConfig:
		logger.Debugf("incompatible config")
		fmt.Println("It looks like you ran the program with an incompatible configuration.\n" +
			"Log in so that the program can create a valid configuration file.")
		ctx.Namespace = context.Namespace{}
		logger.Debugf("run login")
		if err := login.RunLogin(ctx, login.Flags{
			Username: ctx.Client.Username,
			Password: ctx.Client.Password,
		}); err != nil {
			logger.WithError(err).Errorf("unable to login")
			return err
		}
		logger.Debugf("end login")
	default:
		ctx.Log.WithError(err).Errorf("unable to load config")
		return err
	}
	logger.Debugf("running setup")
	err = clisetup.Setup(ctx)
	if err != nil {
		logger.WithError(err).Errorf("unable to run setup")
	} else {
		logger.Debugf("end setup")
	}
	return err
}

func GetNamespaceByUserfriendlyID(ctx *context.Context, flags *pflag.FlagSet) error {
	var userfriendlyID string
	if flags.Changed("namespace") {
		userfriendlyID, _ = flags.GetString("namespace")
	} else {
		return nil
	}
	nsList, err := ctx.Client.GetNamespaceList()
	if err != nil {
		return err
	}
	ns, ok := nsList.GetByUserFriendlyID(userfriendlyID)
	if !ok {
		return fmt.Errorf("unable to find namespace %q", userfriendlyID)
	}
	ctx.SetNamespace(ns)
	return nil
}

func PreRunFunc(ctx *context.Context) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := PreRun(ctx); err != nil {
			angel.Angel(ctx, err)
			os.Exit(1)
		}
		if err := GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func WithInit(ctx *context.Context, action func(*context.Context) *cobra.Command) *cobra.Command {
	var cmd = action(ctx)
	cmd.PreRun = PreRunFunc(ctx)
	return cmd
}

func ResolveLabel(ctx *context.Context, label string) (namespace.Namespace, error) {
	nsList, err := ctx.Client.GetNamespaceList()
	if err != nil {
		return namespace.Namespace{}, err
	}
	ns, ok := nsList.GetByUserFriendlyID(label)
	if !ok {
		return namespace.Namespace{}, fmt.Errorf("unable to find deployment")
	}
	return ns, nil
}
