package model

// represents port protocol type (TCP or UDP)
//
// swagger:model
type Protocol string

// represents service type
//
// swagger:model
type ServiceType string

const (
	// UDP net protocol
	UDP Protocol = "UDP"
	// TCP net protocol
	TCP Protocol = "TCP"
)

// ServicesList -- model for services list
//
// swagger:model
type ServicesList struct {
	Services []Service `json:"services"`
}

// represents service
//
// swagger:model
type Service struct {
	// required: true
	Name string `json:"name"`
	//creation date in RFC3339 format
	CreatedAt *string  `json:"created_at,omitempty"`
	Deploy    string   `json:"deploy,omitempty"`
	IPs       []string `json:"ips,omitempty"`
	Domain    string   `json:"domain,omitempty"`
	// required: true
	Ports []ServicePort `json:"ports"`
	Owner string        `json:"owner,omitempty"`
}

// represent service port
//
// swagger:model
type ServicePort struct {
	// required: true
	Name string `json:"name"`
	Port *int   `json:"port,omitempty"`
	// required: true
	TargetPort int `json:"target_port"`
	// required: true
	Protocol Protocol `json:"protocol"`
}

// Mask removes information not interesting for users
func (svc *Service) Mask() {
	svc.Owner = ""
}
