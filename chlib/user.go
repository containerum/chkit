package chlib

import (
	jww "github.com/spf13/jwalterweatherman"
)

func UserLogin(client *Client, login, password string, np *jww.Notepad) error {
	if err := client.Login(login, password); err != nil {
		return err
	}
	np.FEEDBACK.Printf("Set token: %v", client.UserConfig.Token)
	// get namespaces and set default namespace
	nsResult, err := GetCmdRequestJson(client, KindNamespaces, "", "")
	if err != nil {
		return err
	}
	client.UserConfig.Namespace = nsResult[0]["data"].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
	np.FEEDBACK.Printf("Chosen namespace: %v", client.UserConfig.Namespace)
	return err
}
