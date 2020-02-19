package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"gitlab.informatika.org/label-1-backend/base/auth/usermgmt"
)

// UserStore implements database operations for User management by user.
type UserStore struct {
	db *pg.DB
}

// NewUserStore returns an UserStore.
func NewUserStore(db *pg.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

//Create user
func (s *UserStore) Create(a *usermgmt.User) (*usermgmt.User, error) {

	err := s.db.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, err
}

// Get an User by ID.
func (s *UserStore) Get(id int) (*usermgmt.User, error) {

	a := usermgmt.User{UserID: id}

	err := s.db.Model(&a).Where("user_id = ?", id).Select()

	if err != nil {
		return nil, err
	}

	return &a, nil
}

//Get User by Cookie
func (s *UserStore) GetByCookie(cookie string) (*usermgmt.User, error) {

	a := usermgmt.User{Cookie: cookie}

	err := s.db.Model(&a).Where("cookie = ?", cookie).Select()

	return &a, err
}

// Update a User.
func (s *UserStore) Update(id int, a *usermgmt.User) (*usermgmt.User, error) {
	a.UserID = id
	err := s.db.Update(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Delete an User.
func (s *UserStore) Delete(id int) (*usermgmt.User, error) {

	model := usermgmt.User{UserID: id}

	fmt.Println(model)

	err := s.db.Delete(&model)

	if err != nil {
		return nil, err
	}

	return &model, nil
}