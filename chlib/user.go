package chlib

import (
	"github.com/containerum/chkit/chlib/dbconfig"
	"github.com/containerum/chkit/helpers"

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
	// get namespaces and set default namespace
	nsResult, err := GetCmdRequestJson(client, KindNamespaces, "", "")
	if err != nil {
		return
	}
	info.Namespace = nsResult[0]["data"].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
	np.FEEDBACK.Printf("Chosen namespace: %v", info.Namespace)
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
