package cmd

import (
	"fmt"
	"regexp"

	"chkit-v2/chlib"

	"chkit-v2/helpers"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

const emailRegex = "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open session and set up token",
	PreRun: func(cmd *cobra.Command, args []string) {
		// Hidden password input works incorrect in Windows
		if helpers.IsWindows() {
			if !cmd.Flag("login").Changed || !cmd.Flag("password").Changed {
				np.FEEDBACK.Println("Login and password must be specified")
				cmd.Usage()
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		isValidMail := regexp.MustCompile(emailRegex)
		var email string
		if !cmd.Flag("login").Changed {
			np.FEEDBACK.Print("Enter your email: ")
			fmt.Scan(&email)
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
			passwordB, _ := gopass.GetPasswdMasked()
			password = string(passwordB)
		} else {
			password = cmd.Flag("password").Value.String()
		}
		token, err := chlib.UserLogin(db, email, password, np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		fmt.Println("Successful login\nToken changed to: ", token)
	},
}

func init() {
	loginCmd.PersistentFlags().StringP("login", "l", "", "User login (email)")
	loginCmd.PersistentFlags().StringP("password", "p", "", "User password")
	RootCmd.AddCommand(loginCmd)
}
