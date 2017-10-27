package dbconfig

import (
	"fmt"

	"github.com/containerum/chkit/helpers"
)

type UserInfo struct {
	Token     string `mapconv:"token"`
	Namespace string `mapconv:"namespace"`
}

const userBucket = "user"

func init() {
	defaultInfo := UserInfo{Token: "", Namespace: "default"}
	initializers[userBucket] = helpers.StructToMap(defaultInfo)
}

func (d *ConfigDB) GetUserInfo() (info UserInfo, err error) {
	m, err := d.readTransactional(userBucket)
	if err != nil {
		return info, fmt.Errorf("user bucket get: %s", err)
	}
	err = helpers.FillStruct(&info, m)
	if err != nil {
		return info, fmt.Errorf("user data fill: %s", err)
	}
	return info, nil
}

func (d *ConfigDB) UpdateUserInfo(info UserInfo) error {
	return d.writeTransactional(helpers.StructToMap(info), userBucket)
}
