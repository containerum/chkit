package fingerpint

import (
	"bytes"
	// nolint:gas
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
	userInfo, err := user.Current()
	if err != nil {
		panic("[chkit-cmd] unable to get userInfo data for fingerpint:\n" + err.Error())
	}
	userData += userInfo.Username
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("[chkit-cmd] unable to get net interfaces:\n" +
			err.Error())
	}
	for _, netInterface := range interfaces {
		if !bytes.Equal(netInterface.HardwareAddr, nil) {
			userData += netInterface.HardwareAddr.String()
		}
	}
	//#nosec
	sum := md5.Sum([]byte(userData))
	return hex.EncodeToString(sum[:])
}
