package model

import "time"

// Volume -- volume representation
// provided by resource-service
// https://ch.pages.containerum.net/api-docs/modules/resource-service/index.html#get-namespace
type Volume struct {
	CreateTime       time.Time `json:"create_time"`
	Deleted          bool      `json:"deleted"`
	TariffID         string    `json:"tariff_id"`
	Label            string    `json:"label"`
	Access           string    `json:"access"`
	AccessChangeTime time.Time `json:"access_change_time"`
	Storage          int       `json:"storage"`
	Replicas         int       `json:"replicas"`
}

// CreateVolume --
type CreateVolume struct {
	TariffID string `json:"tariff-id"`
	Label    string `json:"label"`
}

// ResourceUpdateName -- containes new resource name
type ResourceUpdateName struct {
	Label string `json:"label"`
}

// ResourceUpdateUserAccess -- containes user access data
type ResourceUpdateUserAccess struct {
	Username string `json:"username"`
	Access   string `json:"access,omitempty"`
}
