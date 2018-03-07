package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net"

	"github.com/matishsiao/goInfo"
)

func Fingerprint() string {
	userData := goInfo.GetInfo().String()
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("[chkit-cmd] unable to get net interfaces " +
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
