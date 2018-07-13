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

	if strings.TrimSpace(ctx.GetClient().Username) == "" {
		username, err = readLogin()
		if err != nil {
			return err
		}
		if strings.TrimSpace(username) == "" {
			return ErrInvalidUsername
		}
		ctx.GetClient().Username = username
	}

	if strings.TrimSpace(ctx.GetClient().Password) == "" {
		pass, err = readPassword()
		if err != nil {
			return err
		}
		if strings.TrimSpace(pass) == "" {
			return ErrInvalidPassword
		}
		ctx.GetClient().Password = pass
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
				ctx.Exit(1)
			}
			flags.Namespace, _ = command.Flags().GetString("project")
			if err := RunLogin(ctx, flags); err != nil {
				fmt.Println(err)
				ctx.Exit(1)
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
	ctx.GetClient().Username = flags.Username
	ctx.GetClient().Password = flags.Password
	ctx.Changed = true
	logger.Debugf("start app setup")
	if err := Setup(ctx); err != nil {
		angel.Angel(ctx, err)
		ctx.Exit(1)
	}
	logger.Debugf("end setup")

	switch flags.Namespace {
	case "-":
		GetDefaultNS(ctx, true)
	case "":
		GetDefaultNS(ctx, false)
	default:
		nsList, err := ctx.GetClient().GetNamespaceList()
		logger.Debugf("Getting projects list")
		if err != nil {
			logger.WithError(err).Errorf("unable to get namespace list")
			ferr.Println(err)
			ctx.Exit(1)
		}
		var nsName = flags.Namespace
		ns, ok := nsList.GetByUserFriendlyID(nsName)
		if ok {
			ctx.SetNamespace(context.NamespaceFromModel(ns))
			ctx.Changed = true
		} else {
			ferr.Printf("Project %q not found!\n", nsName)
			GetDefaultNS(ctx, false)
		}
	}
	return nil
}
