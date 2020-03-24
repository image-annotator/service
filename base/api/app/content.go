package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.informatika.org/label-1-backend/base/models"
	// "gitlab.informatika.org/content-1-backend/base/auth/jwt"
)

// The list of error types returned from content resource.
var (
	ErrContentValidation    = errors.New("content validation error")
	ErrContentDuplicateEror = errors.New("ERROR #23505 duplicate key value violates unique constraint \"contents_content_name_key\"")
)

// ContentStore defines database operations for content.
type ContentStore interface {
	Create(a *models.Content) (*models.Content, error)
	GetAll() (*[]models.Content, error)
	Get(id int) (*models.Content, error)
	GetByContentName(contentName string) (*[]models.Content, error)
	Update(id int, a *models.Content) (*models.Content, error)
	Delete(id int) (*models.Content, error)
	GetByExactContentName(contentName string) (*models.Content, error)
}

// ContentResource implements content management handler.
type ContentResource struct {
	Store ContentStore
}

// NewContentResource creates and returns an content resource.
func NewContentResource(store ContentStore) *ContentResource {
	return &ContentResource{
		Store: store,
	}
}

func (rs *ContentResource) router(temp *UserResource) *chi.Mux {

	r := chi.NewRouter()
	authSession := []string{"admin", "labeler", "editor"}

	authSessionmw := temp.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		//CRUD STANDARD
		r.Post("/", rs.create)
		r.Get("/", rs.getByContentName)
		r.Get("/{content_id}", rs.get)
		r.Put("/{content_id}", rs.update)
		r.Delete("/{content_id}", rs.delete)

	})
	return r
}

type contentRequest struct {
	*models.Content
}

func (rs *ContentResource) create(w http.ResponseWriter, r *http.Request) {

	var content models.Content

	err := json.NewDecoder(r.Body).Decode(&content)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	respContent, err := rs.Store.Create(&content)

	if err != nil {
		if err.Error() == ErrContentDuplicateEror.Error() {
			respContent, err = rs.Store.GetByExactContentName(content.ContentName)

			if err != nil {
				render.Render(w, r, ErrRender(err))
				return
			}
		} else {
			render.Render(w, r, ErrRender(err))
			return
		}
	}

	render.Respond(w, r, newGlobalResponse(respContent))
}

func (rs *ContentResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "content_id"))

	respContent, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respContent))
}

func (rs *ContentResource) getByContentName(w http.ResponseWriter, r *http.Request) {

	var respContent *[]models.Content
	var err error

	suggestionValue := r.URL.Query().Get("suggestion")

	if suggestionValue == "" {
		fmt.Println(suggestionValue)
		respContent, err = rs.Store.GetAll()
	} else {
		fmt.Println(suggestionValue)
		respContent, err = rs.Store.GetByContentName(suggestionValue)
	}
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respContent))
}

func (rs *ContentResource) update(w http.ResponseWriter, r *http.Request) {

	var content models.Content

	id, err := strconv.Atoi(chi.URLParam(r, "content_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&content)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	getContent, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if content.ContentName != "" {
		getContent.ContentName = content.ContentName
	}

	respContent, err := rs.Store.Update(id, getContent)

	if err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrContentValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respContent))
}

func (rs *ContentResource) delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "content_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	contentResp, err := rs.Store.Delete(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(contentResp))
}
