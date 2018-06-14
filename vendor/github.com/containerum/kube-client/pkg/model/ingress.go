package model

// SelectedIngressesList -- model for ingresses list from all namespaces
//
// swagger:model
type SelectedIngressesList map[string]IngressesList

// IngressesList -- model for ingresses list
//
// swagger:model
type IngressesList struct {
	Ingress []Ingress `json:"ingresses"`
}

// Ingress -- model for ingress
//
// swagger:model
type Ingress struct {
	// required: true
	Name string `json:"name"`
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at,omitempty"`
	//delete date in RFC3339 format
	DeletedAt string `json:"deleted_at,omitempty"`
	// required: true
	Rules []Rule `json:"rules"`
	Owner string `json:"owner,omitempty"`
}

// Rule -- ingress rule
//
// swagger:model
type Rule struct {
	// required: true
	Host      string  `json:"host"`
	TLSSecret *string `json:"tls_secret,omitempty"`
	// required: true
	Path []Path `json:"path"`
}

// Path -- ingress path
//
// swagger:model
type Path struct {
	// required: true
	Path string `json:"path"`
	// required: true
	ServiceName string `json:"service_name"`
	// required: true
	ServicePort int `json:"service_port"`
}

// Mask removes information not interesting for users
func (ingress *Ingress) Mask() {
	ingress.Owner = ""
}
