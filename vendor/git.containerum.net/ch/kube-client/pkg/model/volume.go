package model

import "time"

// ResourceVolume -- volume representation
// provided by resource-service
// https://ch.pages.containerum.net/api-docs/modules/resource-service/index.html#get-namespace
type ResourceVolume struct {
	ID               *string   `json:"id"`
	CreateTime       time.Time `json:"create_time"`
	Deleted          bool      `json:"deleted"`
	TariffID         string    `json:"tariff_id"`
	Label            string    `json:"label"`
	Access           string    `json:"access"`
	AccessChangeTime time.Time `json:"access_change_time"`
	Limited          *bool     `json:"limited"`
	NewAccess        *string   `json:"new_access"`
	Storage          int       `json:"storage"`
	Replicas         int       `json:"replicas"`
}

type ResourceCreateVolume struct {
	TariffID string `json:"tariff-id"`
	Label    string `json:"label"`
}

type ResourceUpdateName struct {
	Label string `json:"label"`
}

type ResourceUpdateUserAccess struct {
	Username string `json:"username"`
	Access   string `json:"access, omitempty"`
}
