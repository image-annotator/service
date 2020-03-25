package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.informatika.org/label-1-backend/base/models"
	// "gitlab.informatika.org/label-1-backend/base/auth/jwt"
)

var (
	ErrAccessControlValidation = errors.New("access control validation error")
)

// AccessControlStore defines database operations for AccessControl.
type AccessControlStore interface {
	Create(a *models.AccessControl) (*models.AccessControl, error)
	Get(id int) (*models.AccessControl, error)
	Update(id int, a *models.AccessControl) (*models.AccessControl, error)
	Delete(id int) (*models.AccessControl, error)
	GetAll() (*[]models.AccessControl, error)
	Count(id int) (int, error)
}

// AccessControlResource implements AccessControl management handler.
type AccessControlResource struct {
	Store AccessControlStore
}

// NewAccessControlResource creates and returns an AccessControl resource.
func NewAccessControlResource(store AccessControlStore) *AccessControlResource {
	return &AccessControlResource{
		Store: store,
	}
}

func (rs *AccessControlResource) router(config API) *chi.Mux {

	r := chi.NewRouter()
	authSession := []string{"admin", "labeler", "editor"}

	authSessionmw := config.User.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		//CRUD STANDARD
		r.Post("/", config.createAC)
		r.Get("/{image_id}", config.getAC)
		r.Get("/", config.getAllAC)
		r.Post("/", config.createAC)
		r.Get("/requestaccess/{image_id}", config.requestAccess)
		r.Put("/{image_id}", config.updateAC)
		r.Delete("/{image_id}", config.deleteAC)
	})
	return r
}

type accessControlRequest struct {
	*models.AccessControl
}

func (rs *API) createAC(w http.ResponseWriter, r *http.Request) {

	var accessControl models.AccessControl

	err := json.NewDecoder(r.Body).Decode(&accessControl)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	respAC, err := rs.AccessControl.Store.Create(&accessControl)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *API) getAC(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	fmt.Println(id, "AIIIDIIII")

	respAC, err := rs.AccessControl.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *API) getAllAC(w http.ResponseWriter, r *http.Request) {

	respAC, err := rs.AccessControl.Store.GetAll()

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *API) requestAccess(w http.ResponseWriter, r *http.Request) {

	var accessgranted bool
	var model models.AccessControl
	var currentAC *models.AccessControl

	accessgranted = false
	imageid, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	user, err := rs.User.Store.GetByCookie(auth[1])

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	model.ImageID = imageid
	model.Timeout = time.Now().Add(time.Minute * 3)
	model.UserID = user.UserID

	currentAC, err = rs.AccessControl.Store.Get(imageid)

	if err != nil {
		if !(err.Error() == "pg: no rows in result set") {
			render.Render(w, r, ErrRender(err))
			return
		}
		currentAC, err = rs.AccessControl.Store.Create(&model)
	}

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if currentAC.Timeout.Before(time.Now()) {
		//EXPIRED

		_, err := rs.AccessControl.Store.Update(imageid, &model)

		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}

		accessgranted = true

	} else {
		//IN EFFECT
		if currentAC.UserID == model.UserID {

			_, err := rs.AccessControl.Store.Update(imageid, &model)

			if err != nil {
				render.Render(w, r, ErrRender(err))
				return
			}

			accessgranted = true
		} else {
			accessgranted = false
		}
	}

	render.Respond(w, r, newGlobalResponse(accessgranted))
}

func (rs *API) updateAC(w http.ResponseWriter, r *http.Request) {

	var accessControl models.AccessControl

	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&accessControl)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	getAC, err := rs.AccessControl.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if &accessControl.Timeout != nil {
		getAC.Timeout = accessControl.Timeout
	}

	if &accessControl.UserID != nil {
		getAC.UserID = accessControl.UserID
	}

	respAC, err := rs.AccessControl.Store.Update(id, getAC)

	if err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrAccessControlValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *API) deleteAC(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	respAC, err := rs.AccessControl.Store.Delete(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}
