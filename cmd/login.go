package cmd

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func readLogin() string {
	notepad.FEEDBACK.Print("Enter your email: ")
	email, err := bufio.NewReader(os.Stdin).ReadString('\n')
	email = strings.TrimRight(email, "\r\n")
	exitOnErr(err)
	return email
}

func readPassword() string {
	notepad.FEEDBACK.Print("Enter your password: ")
	passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
	exitOnErr(err)
	return string(passwordB)
}
