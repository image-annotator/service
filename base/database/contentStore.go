package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"gitlab.informatika.org/label-1-backend/base/models"
)

// ContentStore implements database operations for Content management by content.
type ContentStore struct {
	db *pg.DB
}

// NewContentStore returns an ContentStore.
func NewContentStore(db *pg.DB) *ContentStore {
	return &ContentStore{
		db: db,
	}

}

//Create content
func (s *ContentStore) Create(a *models.Content) (*models.Content, error) {

	err := s.db.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, err
}

//GetAll Get all Content.
func (s *ContentStore) GetAll() (*[]models.Content, error) {

	var contents []models.Content

	err := s.db.Model(&contents).Select()

	if err != nil {
		return nil, err
	}

	return &contents, nil

}

// Get an Content by ID.
func (s *ContentStore) Get(id int) (*models.Content, error) {

	a := models.Content{LabelContentID: id}

	err := s.db.Model(&a).Where("label_contents_id = ?", id).Select()

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetByContentName by imageID an Content by ID.
func (s *ContentStore) GetByContentName(contentName string) (*[]models.Content, error) {

	var contents []models.Content
	queryString := ("%" + contentName + "%")
	fmt.Println(queryString)
	err := s.db.Model(&contents).Where("content_name LIKE ?", queryString).Select()

	if err != nil {
		return nil, err
	}

	return &contents, nil
}

// GetByExactContentName by imageID an Content by ID.
func (s *ContentStore) GetByExactContentName(contentName string) (*models.Content, error) {

	var contents models.Content

	queryString := ("" + contentName + "")
	fmt.Println(queryString)
	err := s.db.Model(&contents).Where("content_name LIKE ?", queryString).Select()

	if err != nil {
		return nil, err
	}

	return &contents, nil
}

// Update a Content.
func (s *ContentStore) Update(id int, a *models.Content) (*models.Content, error) {
	a.LabelContentID = id
	err := s.db.Update(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Delete an Content.
func (s *ContentStore) Delete(id int) (*models.Content, error) {

	models := models.Content{LabelContentID: id}

	fmt.Println(models)

	delContent, err := s.Get(id)

	if err != nil {
		return nil, err
	}

	err = s.db.Delete(&models)

	if err != nil {
		return nil, err
	}

	return delContent, nil
}
