package cmd

import (
	"fmt"
	"regexp"

	"chkit-v2/chlib"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

const emailRegex = "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open session and set up token",
	Run: func(cmd *cobra.Command, args []string) {
		isValidMail := regexp.MustCompile(emailRegex)
		var email string
		if !cmd.Flag("login").Changed {
			jww.FEEDBACK.Print("Enter your email: ")
			fmt.Scan(&email)
		} else {
			email = cmd.Flag("login").Value.String()
		}
		if !isValidMail.MatchString(email) {
			jww.FEEDBACK.Println("Email is not valid")
			return
		}
		var password string
		if !cmd.Flag("password").Changed {
			jww.FEEDBACK.Print("Enter your password: ")
			passwordB, _ := gopass.GetPasswdMasked()
			password = string(passwordB)
		} else {
			password = cmd.Flag("password").Value.String()
		}
		token, err := chlib.UserLogin(db, email, password)
		if err != nil {
			jww.ERROR.Println(err)
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
