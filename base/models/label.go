package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

//Label type struct model
type Label struct {
	LabelID        int       `sql:"label_id,pk" json:"label_id"`
	ImageID        int       `sql:"image_id" json:"image_id"`
	LabelXCenter   float64   `sql:"label_x_center" json:"label_x_center"`
	LabelYCenter   float64   `sql:"label_y_center" json:"label_y_center"`
	LabelWidth     float64   `sql:"label_width" json:"label_width"`
	LabelHeight    float64   `sql:"label_height" json:"label_height"`
	LabelContentID int       `sql:"label_content_id" json:"label_content_id"`
	CreatedAt      time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt      time.Time `sql:"updated_at" json:"updated_at"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Label) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate hook executed before database update operation.
func (a *Label) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()

	return nil
}

// BeforeDelete hook executed before database delete operation.
func (a *Label) BeforeDelete(db orm.DB) error {
	return nil
}
