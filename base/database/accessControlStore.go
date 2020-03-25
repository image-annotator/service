package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"gitlab.informatika.org/label-1-backend/base/models"
)

// AccessControlStore implements database operations for AccessControl management by user.
type AccessControlStore struct {
	db *pg.DB
}

// NewAccessControlStore returns an AccessControlStore.
func NewAccessControlStore(db *pg.DB) *AccessControlStore {
	return &AccessControlStore{
		db: db,
	}
}

func (s *AccessControlStore) Count(id int) (int, error) {

	count, err := s.db.Model(&models.AccessControl{}).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

//Create accesscontrol
func (s *AccessControlStore) Create(a *models.AccessControl) (*models.AccessControl, error) {

	err := s.db.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, err
}

// Get an AccessControl by ID.
func (s *AccessControlStore) Get(id int) (*models.AccessControl, error) {

	a := models.AccessControl{ImageID: id}

	err := s.db.Model(&a).Where("image_id = ?", id).Select()

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetAll AccessControl.
func (s *AccessControlStore) GetAll() (*[]models.AccessControl, error) {

	var accessControls []models.AccessControl

	err := s.db.Model(&accessControls).Select()

	if err != nil {
		return nil, err
	}

	return &accessControls, nil
}

// Update an AccessControl.
func (s *AccessControlStore) Update(id int, a *models.AccessControl) (*models.AccessControl, error) {
	a.ImageID = id
	err := s.db.Update(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Delete an AccessControl.
func (s *AccessControlStore) Delete(id int) (*models.AccessControl, error) {

	model := models.AccessControl{ImageID: id}

	delAccessControl, err := s.Get(id)

	if err != nil {
		return nil, err
	}

	fmt.Println(model)

	err = s.db.Delete(&model)

	if err != nil {
		return nil, err
	}

	return delAccessControl, nil
}
