package model

// model for secret
//
// swagger:model
type Secret struct {
	// required: true
	Name string `json:"name"`
	//creation date in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
	// required: true
	Data map[string]string `json:"data"`
}
