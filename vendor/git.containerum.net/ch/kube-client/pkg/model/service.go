package model

type Protocol string
type ServiceType string

const (
	UDP Protocol = "UDP"
	TCP Protocol = "TCP"
)

const (
	External ServiceType = "external"
	Internal ServiceType = "internal"
)

type Service struct {
	Name      string      `json:"name" binding:"required"`
	CreatedAt *int64      `json:"created_at,omitempty"`
	Deploy    string      `json:"deploy,omitempty"`
	IP        *[]string   `json:"ip,omitempty"`
	Type      ServiceType `json:"type"`
	Ports     []Port      `json:"ports" binding:"required,dive"`
}

type Port struct {
	Name       string   `json:"name" binding:"required"`
	Port       int      `json:"port" binding:"required"`
	TargetPort *int     `json:"target_port,omitempty"`
	Protocol   Protocol `json:"protocol" binding:"required"`
}
