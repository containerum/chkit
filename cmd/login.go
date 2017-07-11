package cmd

import (
	"fmt"
	"regexp"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

const emailRegex = "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open session and set up token",
	Run: func(cmd *cobra.Command, args []string) {
		isValidMail := regexp.MustCompile(emailRegex)
		jww.FEEDBACK.Print("Enter your email: ")
		var email string
		fmt.Scan(&email)
		if !isValidMail.MatchString(email) {
			jww.FEEDBACK.Println("Email is not valid")
			return
		}
		var password string
		jww.FEEDBACK.Print("Enter your password: ")
		fmt.Print("\033[8m") // Hide input
		fmt.Scan(&password)
		fmt.Print("\033[28m") // Show input

		err := chlib.UserLogin(email, password)
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		info, err := chlib.GetUserInfo()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		fmt.Println("Successful login\nToken changed to: ", info.Token)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
