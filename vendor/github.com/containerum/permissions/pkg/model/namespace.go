package model

import (
	"encoding/json"

	"git.containerum.net/ch/permissions/pkg/errors"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

// Namespace describes namespace
//
// swagger:model
type Namespace struct {
	tableName struct{} `sql:"namespaces"`

	Resource

	RAM            int `sql:"ram,notnull" json:"ram"`
	CPU            int `sql:"cpu,notnull" json:"cpu"`
	MaxExtServices int `sql:"max_ext_services,notnull" json:"max_external_services"`
	MaxIntServices int `sql:"max_int_services,notnull" json:"max_internal_services"`
	MaxTraffic     int `sql:"max_traffic,notnull" json:"max_traffic"`

	Volumes []*VolumeWithPermissions `pg:"fk:ns_id" sql:"-" json:"volumes,omitempty"`
}

func (ns *Namespace) BeforeInsert(db orm.DB) error {
	cnt, err := db.Model(ns).
		Where("owner_user_id = ?owner_user_id").
		Where("label = ?label").
		Where("NOT deleted").
		Count()
	if err != nil {
		return err
	}

	if cnt > 0 {
		return errors.ErrResourceAlreadyExists().AddDetailF("namespace %s already exists", ns.Label)
	}

	return nil
}

func (ns *Namespace) AfterInsert(db orm.DB) error {
	return db.Insert(&Permission{
		ResourceID:         ns.ID,
		UserID:             ns.OwnerUserID,
		ResourceType:       ResourceNamespace,
		InitialAccessLevel: AccessOwner,
		CurrentAccessLevel: AccessOwner,
	})
}

func (ns *Namespace) BeforeUpdate(db orm.DB) error {
	if ns.Deleted {
		cnt, err := db.Model(&Volume{NamespaceID: &ns.ID}).
			Where("ns_id = ?ns_id").
			Where("NOT deleted").
			Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			logrus.Error("trying to delete namespace with volumes")
			return errors.ErrInternal()
		}
	}
	return nil
}

// swagger:ignore
type NamespaceWithPermissions struct {
	Namespace `pg:",override"`

	Permission Permission `pg:"fk:resource_id" sql:"-" json:",inline"`

	Permissions []Permission `pg:"polymorphic:resource_" sql:"-" json:"users"`
}

// NamespaceWithPermissions is a response object for get requests
//
// swagger:model NamespaceWithPermissions
type NamespaceWithPermissionsJSON struct {
	Namespace
	Permission
	Permissions []Permission `json:"users"`
}

// Workaround while json "inline" tag not inlines fields on marshal
func (np NamespaceWithPermissions) MarshalJSON() ([]byte, error) {
	npJSON := NamespaceWithPermissionsJSON{
		Namespace:   np.Namespace,
		Permission:  np.Permission,
		Permissions: np.Permissions,
	}

	return json.Marshal(npJSON)
}

func (np *NamespaceWithPermissions) UnmarshalJSON(b []byte) error {
	var npJSON NamespaceWithPermissionsJSON
	err := json.Unmarshal(b, &npJSON)
	if err != nil {
		return err
	}
	np.Namespace = npJSON.Namespace
	np.Permission = npJSON.Permission
	np.Permissions = npJSON.Permissions
	return nil
}

func (np *NamespaceWithPermissions) Mask() {
	np.Namespace.Mask()
	np.Permission.Mask()
	if np.Namespace.OwnerUserID != np.Permission.UserID {
		np.Permissions = nil
	}
}

// NamespaceAdminCreateRequest contains parameters for creating namespace without billing
//
// swagger:model
type NamespaceAdminCreateRequest struct {
	Label          string `json:"label" binding:"required"`
	CPU            int    `json:"cpu" binding:"required"`
	Memory         int    `json:"memory" binding:"required"`
	MaxExtServices int    `json:"max_ext_services" binding:"required"`
	MaxIntServices int    `json:"max_int_services" binding:"required"`
	MaxTraffic     int    `json:"max_traffic" binding:"required"`
}

// NamespaceAdminResizeRequest contains parameter for resizing namespace without billing
//
// swagger:model
type NamespaceAdminResizeRequest struct {
	CPU            *int `json:"cpu"`
	Memory         *int `json:"memory"`
	MaxExtServices *int `json:"max_ext_services"`
	MaxIntServices *int `json:"max_int_services"`
	MaxTraffic     *int `json:"max_traffic"`
}

// NamespaceCreateRequest contains parameters for creating namespace
//
// swagger:model
type NamespaceCreateRequest struct {
	// swagger:strfmt uuid
	TariffID string `json:"tariff_id" binding:"required,uuid"`

	Label string `json:"label" binding:"required"`
}

// NamespaceRenameRequest contains parameters for renaming namespace
//
// swagger:model
type NamespaceRenameRequest struct {
	Label string `json:"label" binding:"required"`
}

// NamespaceResizeRequest contains parameters for changing namespace quota
//
// swagger:model
type NamespaceResizeRequest struct {
	// swagger:strfmt uuid
	TariffID string `json:"tariff_id" binding:"required,uuid"`
}
