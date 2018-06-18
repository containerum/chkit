package model

// Volume -- volume representation
//
//swagger:model
type Volume struct {
	Name string `json:"name,omitempty"`
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at,omitempty"`
	//delete date in RFC3339 format
	Status      string                     `json:"status,omitempty"`
	DeletedAt   string                     `json:"deleted_at,omitempty"`
	Owner       string                     `json:"owner,omitempty"`
	OwnerLogin  string                     `json:"owner_login,omitempty"`
	Access      AccessLevel                `json:"access,omitempty"`
	TariffID    string                     `json:"tariff_id,omitempty"`
	Capacity    uint                       `json:"capacity,omitempty"`
	StorageName string                     `json:"storage_name,omitempty"` //AKA StorageClass
	AccessMode  PersistentVolumeAccessMode `json:"access_mode,omitempty"`
	Users       []UserAccess               `json:"users,omitempty"`
}

// VolumesList -- model for volumes list
//
// swagger:model
type VolumesList struct {
	Volumes []Volume `json:"volumes"`
}

// CreateVolume --
//swagger:model
type CreateVolume struct {
	TariffID string `json:"tariff_id"`
	Label    string `json:"label"`
}

// ResourceUpdateName -- contains new resource name
//swagger:model
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

// Mask removes information not interesting for users
func (vol *Volume) Mask() {
	vol.Owner = ""
}
