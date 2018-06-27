package prerun

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/setup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

type LoadNamespaceListMode string

const (
	TemporarySetNamespace           LoadNamespaceListMode = ""
	RunNamespaceSelectionAndPersist LoadNamespaceListMode = "run namespace selection and persist"
)

func (mode LoadNamespaceListMode) String() string {
	switch mode {
	case TemporarySetNamespace:
		return "temporary set namespace"
	default:
		return string(mode)
	}
}

type ClientInitMode string

const (
	DoNotAllowSelfSignedTLSCerts ClientInitMode = ""
	AllowSelfSignedTLSCerts      ClientInitMode = "allow self signed certs"
	DoNotInitClient              ClientInitMode = "don't init client"
)

func (mode ClientInitMode) String() string {
	if mode == DoNotAllowSelfSignedTLSCerts {
		return "don't allow self signed TLS certs"
	}
	return string(mode)
}

type Config struct {
	InitClient             ClientInitMode
	RunLoginOnMissingCreds bool
	NamespaceSelection     LoadNamespaceListMode
	Namespace              string
}

func PreRun(ctx *context.Context, optional ...Config) error {
	var logger = ctx.Log.Component("PreRun")
	logger.Debugf("START")
	defer logger.Debugf("END")
	var config = Config{
		InitClient:             DoNotAllowSelfSignedTLSCerts,
		RunLoginOnMissingCreds: false,
		NamespaceSelection:     TemporarySetNamespace,
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

	logger.Debugf("running client init in '%s' mode", config.InitClient)
	switch config.InitClient {
	case DoNotAllowSelfSignedTLSCerts:
		err = setup.Client(ctx, setup.DoNotAlloSelfSignedTLSCerts)
	case AllowSelfSignedTLSCerts:
		err = setup.Client(ctx, setup.AllowSelfSignedTLSCerts)
	case DoNotInitClient:
		return nil
	default:
		panic(fmt.Sprintf("[prerun.PreRun] unreacheable InitClient mode %q", config.InitClient))
	}
	if err != nil {
		return chkitErrors.Fatal(err)
	}

	logger.Debugf("running namespace selection in '%s' mode, namespace=%s", config.NamespaceSelection, config.Namespace)
	switch config.NamespaceSelection {
	case RunNamespaceSelectionAndPersist:
		var nsList, err = ctx.Client.GetNamespaceList()
		if err != nil {
			return chkitErrors.Fatal(err)
		}
		var ns namespace.Namespace
		switch config.Namespace {
		case "-":
			var ok = false
			ns, ok = nsList.Head()
			if !ok {
				return chkitErrors.FatalString("you have no namespaces")
			}
			ctx.SetNamespace(context.NamespaceFromModel(ns))
		case "":
			if nsList.Len() == 0 {
				return chkitErrors.FatalString("you have no namespaces")
			}
			(&activekit.Menu{
				Items: activekit.StringSelector(nsList.OwnersAndLabels(), func(s string) error {
					ns, _ = nsList.GetByUserFriendlyID(s)
					ctx.SetNamespace(context.NamespaceFromModel(ns))
					return nil
				}),
			}).Run()
		default:
			var tokens = str.SplitS(config.Namespace, "/", 2).Map(strings.TrimSpace)
			var owner, label string
			if tokens.Len() == 2 {
				owner, label = tokens[0], tokens[1]
			} else {
				label = tokens[0]
			}
			logger.Debugf("owner=%q label=%q", owner, label)
			if !ctx.GetNamespace().Match(owner, label) || ctx.GetNamespace().IsEmpty() {
				logger.Debugf("getting namespace list")
				var nsList, err = ctx.Client.GetNamespaceList()
				if err != nil {
					return chkitErrors.Fatal(err)
				}
				logger.Debugf("searching namespace %q", tokens.Join("/"))
				var ns, ok = nsList.GetByUserFriendlyID(tokens.Join("/"))
				if !ok {
					return chkitErrors.FatalString("namespace %s not found", tokens.Join("/"))
				}
				logger.Debugf("%v", ns.OwnerAndLabel())
				ctx.SetNamespace(context.NamespaceFromModel(ns))
			}
		}
	case TemporarySetNamespace:
		var tokens = str.SplitS(config.Namespace, "/", 2).Map(strings.TrimSpace)
		if config.Namespace != "" && !ctx.GetNamespace().Match(tokens.GetDefault(0, ""), tokens.GetDefault(1, "")) {
			logger.Debugf("getting namespace list")
			var nsList, err = ctx.Client.GetNamespaceList()
			if err != nil {
				return chkitErrors.Fatal(err)
			}
			var ns, ok = nsList.GetByUserFriendlyID(tokens.Join("/"))
			if !ok {
				return chkitErrors.FatalString("you have no namespaces")
			}
			ctx.SetTemporaryNamespace(ns)
		}
		logger.Debugf("using namespace %q", ctx.GetNamespace())

	default:
		panic(fmt.Sprintf("[prerun.PreRun] invalid NamespaceSelection mode %q", config.NamespaceSelection))
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
		if len(optional) == 0 {
			optional = []Config{{}}
		}
		for i, opt := range optional {
			if opt.Namespace == "" {
				opt.Namespace, _ = cmd.Flags().GetString("namespace")
			}
			optional[i] = opt
		}
		if err := PreRun(ctx, optional...); err != nil {
			panic(err)
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
