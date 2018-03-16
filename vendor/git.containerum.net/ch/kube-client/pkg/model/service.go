package model

// Protocol -- represents port protocol type
type Protocol string

// ServiceType -- represents service type
type ServiceType string

const (
	// UDP net protocol
	UDP Protocol = "UDP"
	// TCP net protocol
	TCP Protocol = "TCP"
)

// Service --
type Service struct {
	Name      string        `json:"name"`
	CreatedAt *int64        `json:"created_at,omitempty"`
	Deploy    string        `json:"deploy,omitempty"`
	IPs       []string      `json:"ips,omitempty"`
	Domain    string        `json:"domain,omitempty"`
	Ports     []ServicePort `json:"ports"`
}

// ServicePort -- represent service port
type ServicePort struct {
	Name       string   `json:"name"`
	Port       *int     `json:"port,omitempty"`
	TargetPort int      `json:"target_port"`
	Protocol   Protocol `json:"protocol"`
}
