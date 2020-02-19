// Package app ties together application resources and handlers.
package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"

	"gitlab.informatika.org/label-1-backend/base/database"
	"gitlab.informatika.org/label-1-backend/base/logging"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
	ctxProfile
)

// API provides application resources and handlers.
type API struct {
	Account *AccountResource
	Profile *ProfileResource
	User    *UserResource
	Image   *ImageResource
}

// NewAPI configures and returns application API.
func NewAPI(db *pg.DB) (*API, error) {
	accountStore := database.NewAccountStore(db)
	account := NewAccountResource(accountStore)

	profileStore := database.NewProfileStore(db)
	profile := NewProfileResource(profileStore)

	userStore := database.NewUserStore(db)
	user := NewUserResource(userStore)

	imageStore := database.NewImageStore(db)
	image := NewImageResource(imageStore)

	api := &API{
		Account: account,
		Profile: profile,
		User:    user,
		Image:   image,
	}
	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/account", a.Account.router())
	r.Mount("/profile", a.Profile.router())
	r.Mount("/user", a.User.router())
	r.Mount("/image", a.Image.router())

	return r
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
