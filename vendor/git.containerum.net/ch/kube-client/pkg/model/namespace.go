package model

import "time"

type Namespace struct {
	Created   int64     `json:"created_at,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Resources Resources `json:"resources" binding:"required"`
}

type Resources struct {
	Hard Resource  `json:"hard" binding:"required"`
	Used *Resource `json:"used,omitempty"`
}

type Resource struct {
	CPU    string `json:"cpu" binding:"required"`
	Memory string `json:"memory" binding:"required"`
}

type UpdateNamespace struct {
	Resources Resources `json:"resources" binding:"required"`
}

// ResourceNamespace -- namespace representation
// provided by resource-service
// https://ch.pages.containerum.net/api-docs/modules/resource-service/index.html#get-namespace
type ResourceNamespace struct {
	CreateTime       time.Time        `json:"create_time"`
	Deleted          bool             `json:"deleted"`
	TariffID         string           `json:"tariff_id"`
	Label            string           `json:"label"`
	Access           string           `json:"access"`
	AccessChangeTime time.Time        `json:"access_change_time"`
	NewAccessLevel   *string          `json:"new_access "`
	RAM              int              `json:"ram"`
	CPU              int              `json:"cpu"`
	MaxExtService    int              `json:"max_ext_service"`
	MaxIntService    int              `json:"max_int_service"`
	MaxTraffic       int              `json:"max_traffic"`
	Volumes          []ResourceVolume `json:"volumes"`
}

type UpdateNamespaceName struct {
	Label string `json:"label"`
}
