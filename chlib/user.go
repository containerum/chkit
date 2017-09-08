package chlib

import (
	"github.com/containerum/chkit.v2/chlib/dbconfig"
	"github.com/containerum/chkit.v2/helpers"

	jww "github.com/spf13/jwalterweatherman"
)

func UserLogin(db *dbconfig.ConfigDB, login, password string, np *jww.Notepad) (token string, err error) {
	info, err := db.GetUserInfo()
	if err != nil {
		return
	}
	client, err := NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
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
