package login

import "github.com/containerum/chkit/pkg/context"

func init() {
	Command.PersistentFlags().StringVarP(&context.GlobalContext.Client.Username, "username", "u", "", "your account login")
	Command.PersistentFlags().StringVarP(&context.GlobalContext.Client.Password, "password", "p", "", "your account password")
}
