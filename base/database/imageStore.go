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

// GetPerPage 20 Image by ID.
func (s *ImageStore) GetPerPage(page int, perpage int) (*[]models.Image, error) {

	var images []models.Image

	err := s.db.Model(&images).Offset((page - 1) * perpage).Limit(perpage).Select()

	if err != nil {
		return nil, err
	}

	return &images, nil
}

// GetAll Image.
func (s *ImageStore) GetAll() (*[]models.Image, int, error) {

	var images []models.Image

	count, err := s.db.Model(&images).SelectAndCount()

	if err != nil {
		return nil, 0, err
	}

	return &images, count, nil
}

//Query image by filename
func (s *ImageStore) GetByFilename(query string, page int, perpage int) (*[]models.Image, error) {

	var images []models.Image

	err := s.db.Model(&images).Where("file_name LIKE ?", "%"+query+"%").Offset((page - 1) * perpage).Limit(perpage).Select()

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

// Label a Image.
func (s *ImageStore) Label(id int, a *models.Image) (*models.Image, error) {
	a.ImageID = id
	a.Labeled = true

	err := s.db.Update(a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Unlabel a Image.
func (s *ImageStore) Unlabel(id int, a *models.Image) (*models.Image, error) {
	a.ImageID = id
	a.Labeled = false

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
