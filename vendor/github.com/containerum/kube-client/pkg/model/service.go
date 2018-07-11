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
	Name string `json:"name" yaml:"name"`
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	//delete date in RFC3339 format
	DeletedAt string   `json:"deleted_at,omitempty" yaml:"deleted_at,omitempty"`
	Deploy    string   `json:"deploy,omitempty" yaml:"deploy,omitempty"`
	IPs       []string `json:"ips,omitempty" yaml:"ips,omitempty"`
	Domain    string   `json:"domain,omitempty" yaml:"domain,omitempty"`
	//Solution ID (only if service is part of solution)
	SolutionID string `json:"solution_id,omitempty" yaml:"solution_id,omitempty"`
	// required: true
	Ports []ServicePort `json:"ports" yaml:"ports"`
	Owner string        `json:"owner,omitempty" yaml:"owner,omitempty"`
}

// represent service port
//
// swagger:model
type ServicePort struct {
	// required: true
	Name string `json:"name" yaml:"name"`
	Port *int   `json:"port,omitempty" yaml:"port,omitempty"`
	// required: true
	TargetPort int `json:"target_port" yaml:"target_port"`
	// required: true
	Protocol Protocol `json:"protocol" yaml:"protocol"`
}

// Mask removes information not interesting for users
func (svc *Service) Mask() {
	svc.Owner = ""
}
