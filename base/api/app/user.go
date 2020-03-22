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

//COOKIE AND PASSCODE GENERATOR

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func generateString(len int) string {
	rand.Seed(time.Now().UnixNano())
	return randomString(len)
}

// The list of error types returned from user resource.
var (
	ErrUserValidation = errors.New("user validation error")
)

// UserStore defines database operations for user.
type UserStore interface {
	Create(a *usermgmt.User) (*usermgmt.User, error)
	Get(id int) (*usermgmt.User, error)
	Update(id int, a *usermgmt.User) (*usermgmt.User, error)
	Delete(id int) (*usermgmt.User, error)
	GetByLogin(a *usermgmt.User) (*usermgmt.User, error)
	GetByCookie(cookie string) (*usermgmt.User, error)
	GetAll() (*[]usermgmt.User, error)
}

// UserResource implements user management handler.
type UserResource struct {
	Store UserStore
}

// NewUserResource creates and returns an user resource.
func NewUserResource(store UserStore) *UserResource {
	return &UserResource{
		Store: store,
	}
}

func (rs *UserResource) router(temp *UserResource) *chi.Mux {

	r := chi.NewRouter()
	authAdmin := []string{"admin"}
	authSession := []string{"admin", "labeler", "editor"}

	authAdminmw := temp.basicAuthFactory(authAdmin)
	authSessionmw := temp.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authAdminmw)
		r.Post("/register", rs.create)
		r.Put("/{user_id}", rs.update)
		r.Delete("/{user_id}", rs.delete)
	})

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		r.Post("/validatesession", rs.getbycookie)
		r.Get("/{user_id}", rs.get)
		r.Get("/", rs.getAll)
	})

	r.Post("/login", rs.login)
	return r
}

type userRequest struct {
	*usermgmt.User
}

func (rs *UserResource) create(w http.ResponseWriter, r *http.Request) {

	var user usermgmt.User

	json.NewDecoder(r.Body).Decode(&user)

	user.Passcode = generateString(8)
	user.Cookie = generateString(40)

	respUser, err := rs.Store.Create(&user)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *UserResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	respUser, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *UserResource) getAll(w http.ResponseWriter, r *http.Request) {

	respUser, err := rs.Store.GetAll()

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *UserResource) login(w http.ResponseWriter, r *http.Request) {

	var user usermgmt.User

	json.NewDecoder(r.Body).Decode(&user)

	respUser, err := rs.Store.GetByLogin(&user)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respUser))
}

func (rs *UserResource) getbycookie(w http.ResponseWriter, r *http.Request) {

	var user usermgmt.User

	json.NewDecoder(r.Body).Decode(&user)

	userResp, err := rs.Store.GetByCookie(user.Cookie)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(userResp))
}

func (rs *UserResource) update(w http.ResponseWriter, r *http.Request) {

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

func (rs *UserResource) delete(w http.ResponseWriter, r *http.Request) {

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
