package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"gitlab.informatika.org/label-1-backend/base/models"
)

// ImageStore implements database operations for Image management by user.
type ImageStore struct {
	db *pg.DB
}

// NewImageStore returns an ImageStore.
func NewImageStore(db *pg.DB) *ImageStore {
	return &ImageStore{
		db: db,
	}
}

//Create user
func (s *ImageStore) Create(a *models.Image) (*models.Image, error) {

	err := s.db.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, err
}

// Get an Image by ID.
func (s *ImageStore) Get(id int) (*models.Image, error) {

	a := models.Image{ImageID: id}

	err := s.db.Model(&a).Where("image_id = ?", id).Select()

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetAll Image.
func (s *ImageStore) GetAll() (*[]models.Image, error) {

	var images []models.Image

	err := s.db.Model(&images).Select()

	if err != nil {
		return nil, err
	}

	return &images, nil
}

//Query image by filename
func (s *ImageStore) GetByFilename(query string) (*[]models.Image, error) {

	var images []models.Image

	err := s.db.Model(&images).Where("file_name LIKE ?", "%"+query+"%").Select()

	if err != nil {
		return nil, err
	}

	return &images, nil
}

// Update a Image.
func (s *ImageStore) Update(id int, a *models.Image) (*models.Image, error) {
	a.ImageID = id
	err := s.db.Update(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Delete an Image.
func (s *ImageStore) Delete(id int) (*models.Image, error) {

	model := models.Image{ImageID: id}

	delImage, err := s.Get(id)

	if err != nil {
		return nil, err
	}

	fmt.Println(model)

	err = s.db.Delete(&model)

	if err != nil {
		return nil, err
	}

	return delImage, nil
}
