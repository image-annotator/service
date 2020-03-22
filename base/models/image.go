package models

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

//Image type struct
type Image struct {
	ImageID   int       `sql:"image_id,pk" 	 json:"image_id"`
	Filename  string    `sql:"file_name" 	 json:"file_name"`
	ImagePath string    `sql:"image_path" 	 json:"image_path"`
	Labeled   bool      `sql:"labeled"		 json:"labeled"`
	CreatedAt time.Time `sql:"created_at"	 json:"created_at,omitempty"`
	UpdatedAt time.Time `sql:"updated_at"	 json:"updated_at,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Image) BeforeInsertE(db orm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *Image) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Image) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Image struct and returns validation errors.
func (a *Image) Validate() error {
	a.ImagePath = strings.TrimSpace(a.ImagePath)

	return validation.ValidateStruct(a,
		validation.Field(&a.ImagePath, validation.Required, is.ASCII),
	)
}
