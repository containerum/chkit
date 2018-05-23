package model

//go:generate swagger generate spec -m -o ../../swagger.json

// ConfigMapData -- model for config map data
//
// swagger:model
type ConfigMapData map[string]string

// ConfigMap -- model for config map
//
// swagger:model
type ConfigMap struct {
	// required: true
	Name string `json:"name"`
	//creation date in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
	// key-value data
	//
	// required: true
	Data ConfigMapData `json:"data"`
}
