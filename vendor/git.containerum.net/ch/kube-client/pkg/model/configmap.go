package model

// ConfigMap --
type ConfigMap struct {
	Name      string            `json:"name"`
	CreatedAt *string           `json:"created_at,omitempty"`
	Data      map[string]string `json:"data"`
}
