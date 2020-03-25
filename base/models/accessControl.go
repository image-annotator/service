package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

type AccessControl struct {
	tableName struct{} `pg:"access_control"`

	ImageID int       `sql:"image_id,pk" json:"image_id"`
	UserID  int       `sql:"user_id" json:"user_id"`
	Timeout time.Time `sql:"timeout" json:"timeout"`
}

// BeforeInsert hook executed before database insert operation.
func (a *AccessControl) BeforeInsert(db orm.DB) error {
	return nil
}

// BeforeUpdate hook executed before database update operation.
func (a *AccessControl) BeforeUpdate(db orm.DB) error {
	return nil
}

// BeforeDelete hook executed before database delete operation.
func (a *AccessControl) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates AccessControl struct and returns validation errors.
func (a *AccessControl) Validate() error {
	return nil
}
