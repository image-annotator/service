package models

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

//Content type struct model
type Content struct {
	ContentContentID int       `sql:"label_contents_id,pk" json:"label_contents_id"`
	ContentName      string    `sql:"content_name" json:"content_name"`
	CreatedAt        time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt        time.Time `sql:"updated_at" json:"updated_at"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Content) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}

	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *Content) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()

	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Content) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Image struct and returns validation errors.
func (a *Content) Validate() error {
	a.ContentName = strings.TrimSpace(a.ContentName)

	return validation.ValidateStruct(a,
		validation.Field(&a.ContentName, validation.Required, is.ASCII),
	)
}
