package user

import (
	"fmt"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type User kubeModels.User

func UserFromKube(kubeUser kubeModels.User) User {
	return User(kubeUser)
}

func (user User) ToKube() kubeModels.User {
	return kubeModels.User(user)
}

func (user User) String() string {
	return fmt.Sprintf(
		"Login   : %s\n"+
			"Name   : %s %s\n"+
			"Phone  : %s\n"+
			"Company: %s\n",
		user.Login,
		user.Data.FirstName, user.Data.LastName,
		user.Data.Phone,
		user.Data.Company)
}
