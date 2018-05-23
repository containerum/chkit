package model

import (
	"time"

	"git.containerum.net/ch/permissions/pkg/errors"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

// Resource represents common resource information.
//
// swagger:ignore
type Resource struct {
	// swagger:strfmt uuid
	ID string `sql:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`

	CreateTime *time.Time `sql:"create_time,default:now(),notnull" json:"create_time,omitempty"`

	Deleted bool `sql:"deleted,notnull" json:"deleted,omitempty"`

	DeleteTime *time.Time `sql:"delete_time" json:"delete_time,omitempty"`

	// swagger:strfmt uuid
	TariffID *string `sql:"tariff_id,type:uuid" json:"tariff_id,omitempty"`

	// swagger:strfmt uuid
	OwnerUserID string `sql:"owner_user_id,type:uuid,notnull" json:"owner_user_id,omitempty"`

	// swagger:strfmt email
	OwnerUserLogin string `sql:"-" json:"owner_user_login,omitempty"`

	Label string `sql:"label,notnull" json:"label"`
}

func (r *Resource) BeforeDelete(db orm.DB) error {
	// do not allow delete from app
	logrus.Error("record delete not allowed, use update set deleted = true")
	return errors.ErrInternal()
}

func (r *Resource) AfterUpdate(db orm.DB) error {
	if r.Deleted {
		now := time.Now()
		r.DeleteTime = &now
		_, err := db.Model(&Permission{ResourceID: r.ID}).
			Where("resource_id = ?resource_id").
			Delete()
		return err
	}
	return nil
}

func (r *Resource) Mask() {
	r.CreateTime = nil
	r.DeleteTime = nil
	r.OwnerUserID = ""
}
