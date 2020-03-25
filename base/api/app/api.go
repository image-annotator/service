// Package app ties together application resources and handlers.
package app

import (
	"fmt"
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

type globalResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

func newGlobalResponse(a interface{}) *globalResponse {
	resp := &globalResponse{Data: a, Status: "success"}
	return resp
}

// API provides application resources and handlers.
type API struct {
	Account       *AccountResource
	Profile       *ProfileResource
	User          *UserResource
	Image         *ImageResource
	Label         *LabelResource
	Content       *ContentResource
	AccessControl *AccessControlResource
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

	labelStore := database.NewLabelStore(db)
	label := NewLabelResource(labelStore)

	contentStore := database.NewContentStore(db)
	content := NewContentResource(contentStore)

	accessControlStore := database.NewAccessControlStore(db)
	accessControl := NewAccessControlResource(accessControlStore)

	api := &API{
		Account:       account,
		Profile:       profile,
		User:          user,
		Image:         image,
		Label:         label,
		Content:       content,
		AccessControl: accessControl,
	}
	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/account", a.Account.router())
	r.Mount("/profile", a.Profile.router())
	r.Mount("/user", a.User.router(a.User))
	r.Mount("/image", a.Image.router(a.User))
	r.Mount("/label", a.Label.router(*a))
	r.Mount("/content", a.Content.router(a.User))
	r.Mount("/accesscontrol", a.AccessControl.router(a.User))
	r.Get("/uploads/{imagepath}", getImage)

	return r
}

func getImage(w http.ResponseWriter, r *http.Request) {

	imagepath := chi.URLParam(r, "imagepath")
	fmt.Println("PEPEEGAAA")
	fmt.Println(imagepath)
	http.ServeFile(w, r, "uploads/"+imagepath)
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
