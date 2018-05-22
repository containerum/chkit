package model

import (
	"encoding/json"

	"git.containerum.net/ch/permissions/pkg/errors"
	"github.com/go-pg/pg/orm"
)

// Volume describes volume
//
// swagger:model
type Volume struct {
	tableName struct{} `sql:"volumes"`

	Resource

	Active *bool `sql:"active,notnull" json:"active,omitempty"`

	Capacity int `sql:"capacity,notnull" json:"capacity"`

	Replicas int `sql:"replicas,notnull" json:"replicas"`

	// swagger:strfmt uuid
	NamespaceID *string `sql:"ns_id,type:uuid" json:"namespace_id,omitempty"`

	GlusterName string `sql:"gluster_name,notnull" json:"gluster_name,omitempty"`

	// swagger:strfmt uuid
	StorageID string `sql:"storage_id,type:uuid,notnull" json:"storage_id,omitempty"`
}

func (v *Volume) BeforeInsert(db orm.DB) error {
	cnt, err := db.Model(v).
		Where("owner_user_id = ?owner_user_id").
		Where("label = ?label").
		Where("NOT deleted").
		Count()
	if err != nil {
		return err
	}

	if cnt > 0 {
		return errors.ErrResourceAlreadyExists().AddDetailF("volume %s already exists", v.Label)
	}

	_, err = db.Model(&Storage{ID: v.StorageID}).
		WherePK().
		Set("used = used + (?)", v.Capacity).
		Update()

	return err
}

func (v *Volume) AfterUpdate(db orm.DB) error {
	if err := v.Resource.AfterUpdate(db); err != nil {
		return err
	}

	var err error
	if v.Deleted {
		_, err = db.Model(&Storage{ID: v.StorageID}).
			WherePK().
			Set("used = used - ?", v.Capacity).
			Update()
	} else {
		oldCapacityQuery := db.Model(v).Column("capacity").WherePK()
		_, err = db.Model(&Storage{ID: v.StorageID}).
			WherePK().
			Set("used = used - (?) + ?", oldCapacityQuery, v.Capacity).
			Update(v)
	}
	return err
}

func (v *Volume) AfterInsert(db orm.DB) error {
	return db.Insert(&Permission{
		ResourceID:         v.ID,
		UserID:             v.OwnerUserID,
		ResourceType:       ResourceVolume,
		InitialAccessLevel: AccessOwner,
		CurrentAccessLevel: AccessOwner,
	})
}

func (v *Volume) Mask() {
	v.Resource.Mask()
	v.Active = nil
	v.Replicas = 0
	v.NamespaceID = nil
	v.GlusterName = ""
	v.StorageID = ""
}

// swagger:ignore
type VolumeWithPermissions struct {
	Volume `pg:",override"`

	Permission Permission `pg:"fk:resource_id" sql:"-" json:",inline"`

	Permissions []Permission `pg:"polymorphic:resource_" sql:"-" json:"users"`
}

// VolumeWithPermissions is a response object for get requests
//
// swagger:model VolumeWithPermissions
type VolumeWithPermissionsJSON struct {
	Volume
	Permission
	Permissions []Permission `json:"users"`
}

// Workaround while json "inline" tag not inlines fields on marshal
func (vp VolumeWithPermissions) MarshalJSON() ([]byte, error) {
	npJSON := VolumeWithPermissionsJSON{
		Volume:      vp.Volume,
		Permission:  vp.Permission,
		Permissions: vp.Permissions,
	}

	return json.Marshal(npJSON)
}

func (vp *VolumeWithPermissions) UnmarshalJSON(b []byte) error {
	var vpJSON VolumeWithPermissionsJSON
	if err := json.Unmarshal(b, &vpJSON); err != nil {
		return err
	}
	vp.Volume = vpJSON.Volume
	vp.Permissions = vpJSON.Permissions
	vp.Permission = vpJSON.Permission
	return nil
}

func (vp *VolumeWithPermissions) Mask() {
	vp.Volume.Mask()
	vp.Permission.Mask()
	if vp.OwnerUserID != vp.Permission.UserID {
		vp.Permissions = nil
	}
}

// VolumeCreateRequest is a request object for creating volume
//
// swagger:model
type VolumeCreateRequest struct {
	// swagger:strfmt uuid
	TariffID string `json:"tariff_id" binding:"required,uuid"`

	Label string `json:"label" binding:"required"`
}

// VolumeRenameRequest is a request object for renaming volume
//
// swagger:model
type VolumeRenameRequest struct {
	Label string `json:"label" binding:"required"`
}

// VolumeResizeRequest contains parameters for changing volume size
//
// swagger:model
type VolumeResizeRequest struct {
	// swagger:strfmt uuid
	TariffID string `json:"tariff_id" binding:"required,uuid"`
}
