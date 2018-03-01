package model

// Ingress --
type Ingress struct {
	Name      string `json:"name"`
	CreatedAt *int64 `json:"created_at,omitempty"`
	Rules     []Rule `json:"rules"`
}

// Rule --
type Rule struct {
	Host      string  `json:"host"`
	TLSSecret *string `json:"tls_secret,omitempty"`
	Path      []Path  `json:"path"`
}

// Path --
type Path struct {
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	ServicePort int    `json:"service_port"`
}
