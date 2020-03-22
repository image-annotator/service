package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

func (rs *AccessControlResource) router(temp *AccessControlResource) *chi.Mux {

	r := chi.NewRouter()
	authSession := []string{"admin", "labeler", "editor"}

	authSessionmw := temp.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		//CRUD STANDARD
		r.Post("/", rs.create)
		r.Get("/{Image_id}", rs.get)
		r.Get("/", rs.getAll)
		r.Put("/{image_id}", rs.update)
		r.Delete("/{Image_id}", rs.delete)
	})
	return r
}

type accessControlRequest struct {
	*models.AccessControl
}

func (rs *AccessControlResource) create(w http.ResponseWriter, r *http.Request) {

	var accessControl models.AccessControl

	err := json.NewDecoder(r.Body).Decode(&accessControl)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	respAC, err := rs.Store.Create(&accessControl)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *AccessControlResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	respAC, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *AccessControlResource) getAll(w http.ResponseWriter, r *http.Request) {

	respAC, err := rs.Store.GetAll()

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}

func (rs *AccessControlResource) update(w http.ResponseWriter, r *http.Request) {

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

	getAC, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if accessControl.Timeout != "" {
		getAC.Timeout = accessControl.Timeout
	}

	if accessControl.AccountID != "" {
		getAC.AccountID = accessControl.AccountID
	}

	respAC, err := rs.Store.Update(id, getAC)

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

func (rs *AccessControlResource) delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	respAC, err := rs.Store.Delete(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respAC))
}
