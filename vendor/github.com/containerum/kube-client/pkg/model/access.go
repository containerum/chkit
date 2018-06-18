package model

type UserAccess struct {
	Username    string      `json:"username"`
	AccessLevel AccessLevel `json:"access_level"`
}

func (access UserAccess) String() string {
	return access.Username + ":" + access.AccessLevel.String()
}

// ResourceUpdateUserAccess -- contains user access data
//swagger:model
type ResourceUpdateUserAccess struct {
	Username string      `json:"username"`
	Access   AccessLevel `json:"access,omitempty"`
}
