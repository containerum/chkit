package cmd

import (
	"regexp"

	"github.com/containerum/chkit/chlib"

	"strings"

	"syscall"

	"bufio"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const emailRegex = "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open session and set up token",
	Run: func(cmd *cobra.Command, args []string) {
		isValidMail := regexp.MustCompile(emailRegex)
		var email string
		if !cmd.Flag("login").Changed {
			np.FEEDBACK.Print("Enter your email: ")
			var err error
			email, err = bufio.NewReader(os.Stdin).ReadString('\n')
			email = strings.TrimRight(email, "\r\n")
			exitOnErr(err)
		} else {
			email = cmd.Flag("login").Value.String()
		}
		if !isValidMail.MatchString(email) {
			np.FEEDBACK.Println("Email is not valid")
			return
		}
		var password string
		if !cmd.Flag("password").Changed {
			np.FEEDBACK.Print("Enter your password: ")
			passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
			exitOnErr(err)
			password = string(passwordB)
		} else {
			password = cmd.Flag("password").Value.String()
		}
		exitOnErr(chlib.UserLogin(client, strings.ToLower(email), password, np))
		saveUserSettings(*client.UserConfig)
	},
}

func init() {
	loginCmd.PersistentFlags().StringP("login", "l", "", "User login (email)")
	loginCmd.PersistentFlags().StringP("password", "p", "", "User password")
	RootCmd.AddCommand(loginCmd)
}
