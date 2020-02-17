package database

import (
	"github.com/go-pg/pg"
	"gitlab.informatika.org/label-1-backend/base/auth/jwt"
	"gitlab.informatika.org/label-1-backend/base/auth/pwdless"
	"gitlab.informatika.org/label-1-backend/base/models"
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

// Get an User by ID.
func (s *UserStore) Get(id int) (*pwdless.User, error) {
	a := pwdless.User{ID: id}
	err := s.db.Model(&a).
		Where("User.id = ?id").
		Column("User.*", "Token").
		First()
	return &a, err
}

// Update an User.
func (s *UserStore) Update(a *pwdless.User) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		WherePK().
		Update()
	return err
}

// Delete an User.
func (s *UserStore) Delete(a *pwdless.User) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
			Where("User_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(&models.Profile{}).
			Where("User_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		return tx.Delete(a)
	})
	return err
}

// UpdateToken updates a jwt refresh token.
func (s *UserStore) UpdateToken(t *jwt.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		WherePK().
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *UserStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}
