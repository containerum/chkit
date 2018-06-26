package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

type Flags struct {
	Username  string `flag:"username u"`
	Password  string `flag:"password p"`
	Namespace string `flag:"-"`
}

var (
	// ErrUnableToReadPassword -- unable to read password
	ErrUnableToReadPassword chkitErrors.Err = "unable to read password"
	// ErrUnableToReadUsername -- unable to read username
	ErrUnableToReadUsername chkitErrors.Err = "unable to read username"
	// ErrInvalidPassword -- invalid password
	ErrInvalidPassword chkitErrors.Err = "invalid password"
	// ErrInvalidUsername -- invalid username
	ErrInvalidUsername chkitErrors.Err = "invalid username"
)

func InteractiveLogin(ctx *context.Context) error {
	var err error
	var username, pass string

	if strings.TrimSpace(ctx.Client.Username) == "" {
		username, err = readLogin()
		if err != nil {
			return err
		}
		if strings.TrimSpace(username) == "" {
			return ErrInvalidUsername
		}
		ctx.Client.Username = username
	}

	if strings.TrimSpace(ctx.Client.Password) == "" {
		pass, err = readPassword()
		if err != nil {
			return err
		}
		if strings.TrimSpace(pass) == "" {
			return ErrInvalidPassword
		}
		ctx.Client.Password = pass
	}
	return nil
}

func readLogin() (string, error) {
	fmt.Print("Enter your email: ")
	email, err := bufio.NewReader(os.Stdin).ReadString('\n')
	email = strings.TrimRight(email, "\r\n")
	if err != nil {
		return "", ErrUnableToReadUsername.Wrap(err)
	}
	return email, nil
}

func readPassword() (string, error) {
	fmt.Print("Enter your password: ")
	passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", ErrUnableToReadPassword.Wrap(err)
	}
	fmt.Println("")
	return string(passwordB), nil
}

func Login(ctx *context.Context) *cobra.Command {
	var flags Flags
	command := &cobra.Command{
		Use:   "login",
		Short: "Login to system",
		Run: func(command *cobra.Command, args []string) {
			if err := SetupLogs(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			flags.Namespace, _ = command.Flags().GetString("namespace")
			if err := RunLogin(ctx, flags); err != nil {
				ferr.Println(err)
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
		GetDefaultNS(ctx, true)
	case "":
		GetDefaultNS(ctx, false)
	default:
		nsList, err := ctx.Client.GetNamespaceList()
		logger.Debugf("Getting namespace list")
		if err != nil {
			logger.WithError(err).Errorf("unable to get namespace lsit")
			ferr.Println(err)
			os.Exit(1)
		}
		var nsName = flags.Namespace
		ns, ok := nsList.GetByUserFriendlyID(nsName)
		if ok {
			ctx.SetNamespace(ns)
			ctx.Changed = true
		} else {
			ferr.Printf("Namespace %q not found!\n", nsName)
			GetDefaultNS(ctx, false)
		}
	}
	return nil
}
