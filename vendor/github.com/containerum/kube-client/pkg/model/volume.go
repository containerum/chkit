package model

import "time"

// Volume -- volume representation
//
//swagger:model
type Volume struct {
	ID               string    `json:"id"`
	CreateTime       time.Time `json:"create_time"`
	Label            string    `json:"label"`
	Access           string    `json:"access"`
	AccessChangeTime time.Time `json:"access_change_time"`
	Capacity         int       `json:"capacity"`
	Replicas         int       `json:"replicas"`
}

// CreateVolume --
//swagger:ignore
type CreateVolume struct {
	TariffID string `json:"tariff-id"`
	Label    string `json:"label"`
}

// ResourceUpdateName -- contains new resource name
//swagger:ignore
type ResourceUpdateName struct {
	Label string `json:"label"`
}

type PersistentVolumeAccessMode string

const (
	// can be mounted read/write mode to exactly 1 host
	ReadWriteOnce PersistentVolumeAccessMode = "ReadWriteOnce"
	// can be mounted in read-only mode to many hosts
	ReadOnlyMany PersistentVolumeAccessMode = "ReadOnlyMany"
	// can be mounted in read/write mode to many hosts
	ReadWriteMany PersistentVolumeAccessMode = "ReadWriteMany"
)

// PersistentVolumeClaim -- persistent volume claim representation
//
//swagger:model
type PersistentVolumeClaim struct {
	// required: true
	Name string `json:"name"`
	//creation date in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
	// required: true
	StorageClass string `json:"storage_class"`
	// required: true
	AccessMode PersistentVolumeAccessMode `json:"access_mode"`
	// required: true
	Size uint `json:"size"`
}
