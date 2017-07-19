package chlib

import (
	"chkit-v2/chlib/dbconfig"
	"chkit-v2/helpers"
)

func UserLogin(db *dbconfig.ConfigDB, login, password string) (token string, err error) {
	info, err := db.GetUserInfo()
	if err != nil {
		return
	}
	client, err := NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4())
	if err != nil {
		return
	}
	token, err = client.Login(login, password)
	if err != nil {
		return
	}
	info.Token = token
	err = db.UpdateUserInfo(info)
	return
}

func UserLogout(db *dbconfig.ConfigDB) error {
	info, err := db.GetUserInfo()
	if err != nil {
		return err
	}
	info.Token = ""
	err = db.UpdateUserInfo(info)
	return err
}
