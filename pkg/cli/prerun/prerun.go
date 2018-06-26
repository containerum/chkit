package prerun

import (
	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/setup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

type LoadNamespaceListMode string

const (
	SelectFirstNamespace  LoadNamespaceListMode = ""
	RunNamespaceSelection LoadNamespaceListMode = "run namespace selection"
	DoNotLoadNamespaces   LoadNamespaceListMode = "don't load namespaces"
)

type ClientInitMode string

const (
	DoNotAllowSelfSignedTLSCerts ClientInitMode = ""
	AllowSelfSignedTLSCerts      ClientInitMode = "allow self signed certs"
	DoNotInitClient                             = "don't init client"
)

type Config struct {
	InitClient             ClientInitMode
	RunLoginOnMissingCreds bool
	NamespaceSelection     LoadNamespaceListMode
}

func PreRun(ctx *context.Context, optional ...Config) error {
	var logger = ctx.Log.Component("PreRun")
	logger.Debugf("START")
	defer logger.Debugf("END")
	var config = Config{
		InitClient:             DoNotAllowSelfSignedTLSCerts,
		RunLoginOnMissingCreds: false,
		NamespaceSelection:     SelectFirstNamespace,
	}
	for _, c := range optional {
		config = c
	}
	logger.StructFields(config)
	setup.SetupLogs(ctx)

	logger.Debugf("syncing config")
	err := configuration.SyncConfig(ctx)
	switch err {
	case nil:
		// pass
	case configuration.ErrIncompatibleConfig:
		logger.WithError(err).Errorf("incompatible config")
		if config.RunLoginOnMissingCreds {
			fmt.Println("It looks like you ran the program with an incompatible configuration.\n" +
				"Run 'chkit login' to create a valid configuration file.")
			ctx.SetNamespace(context.Namespace{})
			logger.Debugf("run login")
			if err := setup.RunLogin(ctx, setup.Flags{
				Username:  ctx.Client.Username,
				Password:  ctx.Client.Password,
				Namespace: "",
			}); err != nil {
				logger.WithError(err).Errorf("unable to login")
				return err
			}
			logger.Debugf("end login")
		} else {
			ferr.Println(err)
			ctx.Exit(1)
		}
	default:
		ctx.Log.WithError(err).Errorf("unable to load config")
		return err
	}
	logger.Debugf("running setup")
	defer logger.Debugf("end setup")

	switch config.InitClient {
	case DoNotAllowSelfSignedTLSCerts:
		return setup.Client(ctx, setup.DoNotAlloSelfSignedTLSCerts)
	case AllowSelfSignedTLSCerts:
		return setup.Client(ctx, setup.AllowSelfSignedTLSCerts)
	case DoNotInitClient:
		return nil
	}
	return nil
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
	ctx.SetTemporaryNamespace(ns)
	return nil
}

func PreRunFunc(ctx *context.Context, optional ...Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := PreRun(ctx, optional...); err != nil {
			angel.Angel(ctx, err)
			ctx.Exit(1)
		}
		if err := GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
			ferr.Println(err)
			ctx.Exit(1)
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
		return namespace.Namespace{}, fmt.Errorf("unable to find deployment %q", ns)
	}
	return ns, nil
}
