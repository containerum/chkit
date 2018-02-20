package model

type Ingress struct {
	Name      string  `json:"name" binding:"required"`
	CreatedAt *int64  `json:"created_at,omitempty"`
	TLSSecret *string `json:"tls_secret,omitempty"`
	Rule      Rule    `json:"rule" binding:"required"`
}

type Rule struct {
	Host string `json:"host" binding:"required"`
	Path Path   `json:"path" binding:"required"`
}

type Path struct {
	Path        string `json:"path" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	ServicePort int    `json:"service_port" binding:"required"`
}

type ResourceIngress struct {
	Domain    string             `json:"domain"`
	Type      string             `json:"type"`
	CreatedAt *int64             `json:"created_at,omitempty"`
	Service   string             `json:"service"`
	TLS       *ResourceTLSsecret `json:"tls, omitempty"`
}

type ResourceTLSsecret struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}
