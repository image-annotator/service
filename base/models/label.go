package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

type Image struct {
	LabelID        int `sql:"label_id,pk" 				json:"label_id"`
	ImageID        int `sql:"image_id" 					json:"image_id"`
	LabelXCenter   int `sql:"label_x_center" 			json:"label_x_center"`
	LabelYCenter   int `sql:"label_y_center" 			json:"label_y_center"`
	LabelWidth     int `sql:"label_width" 				json:"label_width"`
	LabelHeight    int `sql:"label_height" 				json:"label_height"`
	LabelContentID int `sql:"label_content_id" 			json:"label_content_id"`
	CreatedAt      int `sql:"created_at" 				json:"created_at"`
	UpdatedAt      int `sql:"updated_at" 				json:"updated_at"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Image) BeforeInsert(db orm.DB) error {
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
}

// BeforeDelete hook executed before database delete operation.
func (a *Image) BeforeDelete(db orm.DB) error {
	return nil
}
