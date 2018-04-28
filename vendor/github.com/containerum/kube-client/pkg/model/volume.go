package model

import "time"

// Volume -- volume representation
// provided by resource-service
// https://ch.pages.containerum.net/api-docs/modules/resource-service/index.html#get-namespace
//
//swagger:model
type Volume struct {
	CreateTime       time.Time `json:"create_time"`
	Label            string    `json:"label"`
	Access           string    `json:"access"`
	AccessChangeTime time.Time `json:"access_change_time"`
	Storage          int       `json:"storage"`
	Replicas         int       `json:"replicas"`
}

// CreateVolume --
//swagger:ignore
type CreateVolume struct {
	TariffID string `json:"tariff-id"`
	Label    string `json:"label"`
}

// ResourceUpdateName -- containes new resource name
//swagger:ignore
type ResourceUpdateName struct {
	Label string `json:"label"`
}

// ResourceUpdateUserAccess -- containes user access data
//swagger:ignore
type ResourceUpdateUserAccess struct {
	Username string `json:"username"`
	Access   string `json:"access,omitempty"`
}
