package model

import (
	"bytes"
	"fmt"
	"time"

	"git.containerum.net/ch/permissions/pkg/errors"
	"github.com/go-pg/pg/orm"
)

// swagger:ignore
type AccessLevel string // enum

const (
	AccessNone       AccessLevel = "none"
	AccessRead       AccessLevel = "read"
	AccessReadDelete AccessLevel = "readdelete"
	AccessWrite      AccessLevel = "write"
	AccessOwner      AccessLevel = "owner"
)

func (al *AccessLevel) UnmarshalJSON(b []byte) error {
	str := AccessLevel(bytes.Replace(b, []byte(`"`), []byte(``), -1))
	switch str {
	case AccessNone, AccessRead, AccessReadDelete, AccessWrite, AccessOwner:
		*al = str
	default:
		return fmt.Errorf("invalid acess level %s", b)
	}

	return nil
}

type ResourceType string

const (
	ResourceNamespace ResourceType = "Namespace"
	ResourceVolume    ResourceType = "Volume"
)

// Permission represents information about user permission to resource
//
// swagger:model
type Permission struct {
	tableName struct{} `sql:"permissions"`

	// swagger:strfmt uuid
	ID string `sql:"perm_id,pk,type:uuid,default:uuid_generate_v4()" json:"perm_id,omitempty"`

	ResourceType ResourceType `sql:"resource_type,notnull,unique:unique_user_access" json:"kind,omitempty"` // WARN: custom type here, do not forget create it

	// swagger:strfmt uuid
	ResourceID string `sql:"resource_id,type:uuid,notnull,unique:unique_user_access" json:"resource_id,omitempty"`

	CreateTime *time.Time `sql:"create_time,default:now(),notnull" json:"create_time,omitempty"`

	// swagger:strfmt uuid
	UserID string `sql:"user_id,type:uuid,notnull,unique:unique_user_access" json:"user_id,omitempty"`

	// swagger:strfmt email
	UserLogin string `sql:"-" json:"login,omitempty"`

	InitialAccessLevel AccessLevel `sql:"initial_access_level,type:ACCESS_LEVEL,notnull" json:"access,omitempty"` // WARN: custom type here, do not forget create it

	CurrentAccessLevel AccessLevel `sql:"current_access_level,type:ACCESS_LEVEL,notnull" json:"new_access_level,omitempty"` // WARN: custom type here, do not forget create it

	AccessLevelChangeTime *time.Time `sql:"access_level_change_time,default:now(),notnull" json:"access_level_change_time,omitempty"`
}

func (p *Permission) BeforeInsert(db orm.DB) error {
	if p.InitialAccessLevel == AccessOwner {
		cnt, err := db.Model(p).
			Column("id").
			Where("resource_id = ?", p.ResourceID).
			Where("resource_type = ?", p.ResourceType).
			Where("initial_access_level = ?", AccessOwner).
			Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			return errors.ErrOwnerAlreadyExists()
		}
	}
	if p.CurrentAccessLevel > p.InitialAccessLevel {
		// that`s our error if we will here
		return errors.ErrInternal().AddDetails("initial access level must be greater than current access level")
	}
	return nil
}

func (p *Permission) BeforeUpdate(db orm.DB) error {
	if p.CurrentAccessLevel > p.InitialAccessLevel {
		// that`s our error if we will here
		return errors.ErrInternal().AddDetails("initial access level must be greater than current access level")
	}
	if p.InitialAccessLevel == AccessOwner && p.CurrentAccessLevel != p.InitialAccessLevel { // limiting access
		_, err := db.Model(p).
			Where("resource_id = ?resource_id").
			Set("current_access_level = LEAST(?TableAlias.initial_access_level, ?current_access_level)::ACCESS_LEVEL").
			Update()
		if err != nil {
			return err
		}
	}
	if p.InitialAccessLevel == AccessOwner && p.InitialAccessLevel == p.CurrentAccessLevel { // un-limiting access
		_, err := db.Model(p).
			Where("resource_id = ?resource_id").
			Set("current_access_level = initial_access_level").
			Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Permission) Mask() {
	p.ID = ""
	p.ResourceType = "" // will be already known though
	p.ResourceID = ""
	p.CreateTime = nil
	p.UserID = ""
	p.InitialAccessLevel = p.CurrentAccessLevel
	p.AccessLevelChangeTime = nil
	p.CurrentAccessLevel = ""
}

// SetUserAccessRequest is a request object for setting user accesses
//
// swagger:model SetResourcesAccessesRequest
type SetUserAccessesRequest struct {
	Access AccessLevel `json:"access"`
}

// SetUserAccessRequest is a request object for setting access to resource for user
//
// swagger:model SetResourceAccessRequest
type SetUserAccessRequest struct {
	// swagger:strfmt email
	UserName string      `json:"username"`
	Access   AccessLevel `json:"access"`
}

// DeleteUserAccessRequest is a request object for deleting access to resource for user
//
// swagger:model DeleteResourceAccessRequest
type DeleteUserAccessRequest struct {
	// swagger:strfmt email
	UserName string `json:"username"`
}
