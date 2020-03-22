package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"

	// "gitlab.informatika.org/label-1-backend/base/auth/jwt"

	"math/rand"
	"time"

	"gitlab.informatika.org/label-1-backend/base/auth/usermgmt"
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
	authAdmin := []string{"admin"}
	authSession := []string{"admin", "labeler", "editor"}

	authAdminmw := temp.basicAuthFactory(authAdmin)
	authSessionmw := temp.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		r.Post("/createaccess", rs.create)
		r.Get("/{Image_id}", rs.get)
		r.Get("/", rs.getAll)
		r.Put("/update", rs.update)
		r.Delete("/{Image_id}", rs.delete)
	})

	r.Post("/login", rs.login)
	return r
}

func (rs *AccessControlResource) create(w http.ResponseWriter, r *http.Request) {

	var ac models.AccessControl

	json.NewDecoder(r.Body).Decode(&user)

	user.Passcode = generateString(8)
	user.Cookie = generateString(40)

	respAccessControl, err := rs.Store.Create(&user)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *AccessControlResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	respUser, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *AccessControlResource) getAll(w http.ResponseWriter, r *http.Request) {

	respUser, err := rs.Store.GetAll()

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *AccessControlResource) login(w http.ResponseWriter, r *http.Request) {

	var user usermgmt.User

	json.NewDecoder(r.Body).Decode(&user)

	respUser, err := rs.Store.GetByLogin(&user)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *AccessControlResource) getbycookie(w http.ResponseWriter, r *http.Request) {

	var user usermgmt.User

	json.NewDecoder(r.Body).Decode(&user)

	userResp, err := rs.Store.GetByCookie(user.Cookie)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(userResp))
}

func (rs *AccessControlResource) update(w http.ResponseWriter, r *http.Request) {

	var user usermgmt.User

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	getUser, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if user.UserRole != "" {
		getUser.UserRole = user.UserRole
	}

	if user.Username != "" {
		getUser.Username = user.Username
	}

	if user.Passcode != "" {
		getUser.Passcode = user.Passcode
	}

	respUser, err := rs.Store.Update(id, getUser)

	if err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrUserValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *AccessControlResource) delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	userResp, err := rs.Store.Delete(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(userResp))
}
