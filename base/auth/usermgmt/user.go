package usermgmt

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

// User represents an authenticated application user
type User struct {
	tableName struct{} `sql:"users" pg:"users"`

	UserID    int       `sql:"user_id,pk" json:"user_id"`
	CreatedAt time.Time `sql:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `sql:"updated_at" json:"updated_at,omitempty"`
	Username  string    `sql:"username,unique" json:"username,omitempty"`
	Cookie    string    `sql:"cookie" json:"cookie,omitempty"`
	Passcode  string    `sql:"passcode" json:"passcode,omitempty"`
	UserRole  string    `sql:"user_role" json:"user_role,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (a *User) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *User) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *User) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates User struct and returns validation errors.
func (a *User) Validate() error {
	a.Username = strings.TrimSpace(a.Username)

	return validation.ValidateStruct(a,
		validation.Field(&a.Username, validation.Required, is.ASCII),
	)
}
