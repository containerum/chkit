package update

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func gracefulRestart() {
	args := make([]string, 0)
	if len(args) > 1 {
		args = os.Args[1:]
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logrus.WithError(err).Error("graceful restart failed")
	}
}
