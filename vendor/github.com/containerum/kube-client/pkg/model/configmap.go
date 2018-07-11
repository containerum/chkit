package model

// ConfigMapData -- model for config map data
//
// swagger:model
type ConfigMapData map[string]string

// ConfigMap -- model for config map
//
// swagger:model
type ConfigMap struct {
	// required: true
	Name string `json:"name" yaml:"name"`
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	//delete date in RFC3339 format
	DeletedAt string `json:"deleted_at,omitempty" yaml:"deleted_at,omitempty"`
	// key-value data
	//
	// required: true
	Data  ConfigMapData `json:"data" yaml:"data"`
	Owner string        `json:"owner,omitempty" yaml:"owner,omitempty"`
}

// SelectedConfigMapsList -- model for config maps list from all namespaces
//
// swagger:model
type SelectedConfigMapsList map[string]ConfigMapsList

// ConfigMapsList -- model for config maps list
//
// swagger:model
type ConfigMapsList struct {
	ConfigMaps []ConfigMap `json:"configmaps"`
}

// Mask removes information not interesting for users
func (cm *ConfigMap) Mask() {
	cm.Owner = ""
}
