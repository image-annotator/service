package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"gitlab.informatika.org/label-1-backend/base/models"
)

// LabelStore implements database operations for Label management by label.
type LabelStore struct {
	db *pg.DB
}

// NewLabelStore returns an LabelStore.
func NewLabelStore(db *pg.DB) *LabelStore {
	return &LabelStore{
		db: db,
	}
}

//Create label
func (s *LabelStore) Create(a *models.Label) (*models.Label, error) {

	err := s.db.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, err
}

// Get all Label.
func (s *LabelStore) GetAll() (*[]models.Label, error) {

	var labels []models.Label

	err := s.db.Model(&labels).Select()

	if err != nil {
		return nil, err
	}

	return &labels, nil

}

// Get an Label by ID.
func (s *LabelStore) Get(id int) (*models.Label, error) {

	a := models.Label{LabelID: id}

	err := s.db.Model(&a).Where("label_id = ?", id).Select()

	if err != nil {
		return nil, err
	}

	return &a, nil
}

//Get Label by Cookie
func (s *LabelStore) GetByCookie(cookie string) (*models.Label, error) {

	a := models.Label{Cookie: cookie}

	err := s.db.Model(&a).Where("cookie = ?", cookie).Select()

	return &a, err
}

//Get Label by Labelname and Passcode
func (s *LabelStore) GetByLogin(a *models.Label) (*models.Label, error) {

	models := new(models.Label)

	err := s.db.Model(models).Where("labelname = ?", a.Labelname).Where("passcode = ?", a.Passcode).Select()
	if err != nil {
		return nil, err
	}

	return models, err
}

// Update a Label.
func (s *LabelStore) Update(id int, a *models.Label) (*models.Label, error) {
	a.LabelID = id
	err := s.db.Update(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Delete an Label.
func (s *LabelStore) Delete(id int) (*models.Label, error) {

	models := models.Label{LabelID: id}

	fmt.Println(models)

	delLabel, err := s.Get(id)

	if err != nil {
		return nil, err
	}

	err = s.db.Delete(&models)

	if err != nil {
		return nil, err
	}

	return delLabel, nil
}
