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

// The list of error types returned from label resource.
var (
	ErrLabelValidation = errors.New("label validation error")
)

// LabelStore defines database operations for label.
type LabelStore interface {
	Create(a *models.Label) (*models.Label, error)
	GetAll() (*[]models.Label, error)
	Get(id int) (*models.Label, error)
	GetByImageID(id int) (*[]models.Label, error)
	GetByContentID(id int) (*[]models.Label, error)
	Update(id int, a *models.Label) (*models.Label, error)
	Delete(id int) (*models.Label, error)
}

// LabelResource implements label management handler.
type LabelResource struct {
	Store LabelStore
}

// NewLabelResource creates and returns an label resource.
func NewLabelResource(store LabelStore) *LabelResource {
	return &LabelResource{
		Store: store,
	}
}

func (rs *LabelResource) router(temp *UserResource) *chi.Mux {

	r := chi.NewRouter()
	authSession := []string{"admin", "labeler", "editor"}

	authSessionmw := temp.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		//STANDARD CRUD
		r.Post("/", rs.create)
		r.Post("/many", rs.createMany)
		r.Get("/", rs.getAll)
		r.Get("/{label_id}", rs.get)
		r.Put("/{label_id}", rs.update)
		r.Delete("/{label_id}", rs.delete)

		//CUSTOM API
		r.Get("/contentquery/{content_id}", rs.getByContentID)
		r.Get("/imagequery/{image_id}", rs.getByImageID)
	})
	return r
}

type payload struct {
	message string
}

type labelRequest struct {
	*models.Label
}

func (rs *LabelResource) create(w http.ResponseWriter, r *http.Request) {

	var label models.Label

	err := json.NewDecoder(r.Body).Decode(&label)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	respLabel, err := rs.Store.Create(&label)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respLabel))
}

func (rs *LabelResource) createMany(w http.ResponseWriter, r *http.Request) {

	var labels []models.Label
	var returnLabels []models.Label

	json.NewDecoder(r.Body).Decode(&labels)

	for _, label := range labels {

		currentTarget, err := rs.Store.Create(&label)

		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}

		returnLabels = append(returnLabels, *currentTarget)
	}

	render.Respond(w, r, newGlobalResponse(returnLabels))
}

func (rs *LabelResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "label_id"))

	respLabel, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respLabel))
}

func (rs *LabelResource) getByContentID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "content_id"))

	respLabel, err := rs.Store.GetByContentID(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respLabel))
}

func (rs *LabelResource) getByImageID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	respLabel, err := rs.Store.GetByImageID(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respLabel))
}

func (rs *LabelResource) getAll(w http.ResponseWriter, r *http.Request) {

	respLabel, err := rs.Store.GetAll()

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respLabel))
}

func (rs *LabelResource) update(w http.ResponseWriter, r *http.Request) {

	var label models.Label

	id, err := strconv.Atoi(chi.URLParam(r, "label_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&label)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	getUser, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	label.CreatedAt = getUser.CreatedAt

	respLabel, err := rs.Store.Update(id, &label)

	if err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrLabelValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respLabel))
}

func (rs *LabelResource) delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "label_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	labelResp, err := rs.Store.Delete(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(labelResp))
}
