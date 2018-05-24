package model

type Permission struct {
	TariffID      string       `json:"tariff_id"`
	Label         string       `json:"label"`
	Access        string       `json:"access"`
	RAM           int          `json:"ram"`
	CPU           int          `json:"cpu"`
	MaxExtService int          `json:"max_ext_service"`
	MaxIntService int          `json:"max_int_service"`
	Users         []UserAccess `json:"users"`
}

func (perm Permission) HasAccess(username string) bool {
	for _, user := range perm.Users {
		if user.Username == username {
			return true
		}
	}
	return false
}

func (perm Permission) GetAccess(username string) (UserAccess, bool) {
	for _, user := range perm.Users {
		if user.Username == username {
			return user, true
		}
	}
	return UserAccess{
		Username:    username,
		AccessLevel: None,
	}, false
}

type UserAccess struct {
	Username    string      `json:"username"`
	AccessLevel AccessLevel `json:"access_level"`
}

func (access UserAccess) String() string {
	return access.Username + ":" + access.AccessLevel.String()
}

// ResourceUpdateUserAccess -- contains user access data
//swagger:ignore
type ResourceUpdateUserAccess struct {
	Username string `json:"username"`
	Access   string `json:"access,omitempty"`
}
