package login

import (
	. "github.com/containerum/chkit/pkg/context"
)

func init() {
	Command.PersistentFlags().StringVarP(&Context.Client.Username, "username", "u", "", "your account login")
	Command.PersistentFlags().StringVarP(&Context.Client.Password, "password", "p", "", "your account password")
}
