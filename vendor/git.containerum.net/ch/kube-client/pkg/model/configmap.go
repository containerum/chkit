package model

//go:generate swagger generate spec -m -o ../../kube-client-swagger.json

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
	Data map[string]string `json:"data"`
}
