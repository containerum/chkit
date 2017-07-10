package chlib

import (
	"fmt"

	"github.com/kfeofantov/chkit-v2/helpers"
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

func GetUserInfo() (info UserInfo, err error) {
	m, err := readFromBucket(userBucket)
	if err != nil {
		return info, fmt.Errorf("user bucket get: ", err)
	}
	err = helpers.FillStruct(&info, m)
	if err != nil {
		return info, fmt.Errorf("user data fill: ", err)
	}
	return info, nil
}

func UpdateUserInfo(info UserInfo) error {
	return pushToBucket(userBucket, helpers.StructToMap(info))
}
