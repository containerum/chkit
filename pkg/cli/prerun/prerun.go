package prerun

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

func PreRun(ctx *context.Context) error {
	clisetup.SetupLogs(ctx)
	logrus.Debugf("loading config")
	if err := configuration.SyncConfig(ctx); err != nil {
		logrus.WithError(err).Errorf("unable to load config")
		return err
	}
	return clisetup.Setup(ctx)
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
		return fmt.Errorf("unable to find deployment")
	}
	ctx.SetNamespace(ns)
	return nil
}

func WithInit(ctx *context.Context, action func(*context.Context) *cobra.Command) *cobra.Command {
	var cmd = action(ctx)
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if err := PreRun(ctx); err != nil {
			angel.Angel(ctx, err)
			os.Exit(1)
		}
		if err := GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
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
