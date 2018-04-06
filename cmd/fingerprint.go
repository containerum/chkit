package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net"
	"os/user"
	"runtime"

	"github.com/matishsiao/goInfo"
)

// Fingerprint -- generates user computer fingerprint.
// Panics on error
func Fingerprint() string {
	userData := goInfo.GetInfo().String() +
		runtime.GOOS +
		runtime.GOARCH +
		runtime.Version() +
		runtime.Compiler
	user, err := user.Current()
	if err != nil {
		panic("[chkit-cmd] unable to get user data for fingerpint:\n" + err.Error())
	}
	userData += user.Username
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("[chkit-cmd] unable to get net interfaces:\n" +
			err.Error())
	}
	for _, netInterface := range interfaces {
		if bytes.Compare(netInterface.HardwareAddr, nil) != 0 {
			userData += netInterface.HardwareAddr.String()
		}
	}
	sum := md5.Sum([]byte(userData))
	return hex.EncodeToString(sum[:])
}
