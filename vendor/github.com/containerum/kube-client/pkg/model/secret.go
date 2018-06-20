package model

// SecretsList -- model for secrets list
//
// swagger:model
type SecretsList struct {
	Secrets []Secret `json:"secrets"`
}

// model for secret
//
// swagger:model
type Secret struct {
	// required: true
	Name string `json:"name"`
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at,omitempty"`
	//delete date in RFC3339 format
	DeletedAt string `json:"deleted_at,omitempty"`
	// required: true
	Data  map[string]string `json:"data"`
	Owner string            `json:"owner,omitempty"`
}

// Mask removes information not interesting for users
func (secret *Secret) Mask() {
	secret.Owner = ""
}
