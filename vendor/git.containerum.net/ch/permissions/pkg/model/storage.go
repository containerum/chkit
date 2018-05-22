package model

import (
	"git.containerum.net/ch/permissions/pkg/errors"
	"github.com/go-pg/pg/orm"
)

// Storage describes volumes storage
//
// swagger:model
type Storage struct {
	tableName struct{} `sql:"storages"`

	// swagger:strfmt uuid
	ID string `sql:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`

	Name string `sql:"name,notnull,unique" json:"name"`

	Size int `sql:"size,notnull" json:"size"`

	Used int `sql:"used,notnull" json:"used"`

	Replicas int `sql:"replicas,notnull" json:"replicas"`

	IPs []string `sql:"ips,notnull,type:inet[],array" json:"ips"`

	Volumes []*Volume `pg:"fk:storage_id" sql:"-" json:"volumes"`
}

func (s *Storage) BeforeInsert(db orm.DB) error {
	cnt, err := db.Model(s).Where("name = ?name").Count()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.ErrResourceAlreadyExists().AddDetailF("storage %s already exists", s.Name)
	}
	return nil
}

func (s *Storage) BeforeUpdate(db orm.DB) error {
	if s.Size < s.Used {
		return errors.ErrQuotaExceeded().AddDetailF("storage quota exceeded")
	}
	return nil
}

func (s *Storage) BeforeDelete(db orm.DB) error {
	cnt, err := db.Model(&Volume{StorageID: s.ID}).
		Where("storage_id = ?storage_id").
		Where("NOT deleted").
		Count()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.ErrStorageDelete()
	}
	return nil
}

// UpdateStorageRequest represents request object for updating storage
//
// swagger:model
type UpdateStorageRequest struct {
	Name     *string  `json:"name,omitempty"`
	Size     *int     `json:"size,omitempty"`
	Replicas *int     `json:"replicas,omitempty"`
	IPs      []string `json:"ips,omitempty"`
}
