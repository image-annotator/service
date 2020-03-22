package models

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

type AccessControl struct {
	ImageID   int       `sql:"image_id,pk" json:"image_id"`
	AccountID int       `json:"-"`
	Timeout time.Time 	`json:"updated_at,omitempty"`
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
