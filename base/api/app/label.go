package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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
		r.Post("/create", rs.create)
		r.Post("/createMany", rs.createMany)
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

	json.NewDecoder(r.Body).Decode(&labels)

	for _, label := range labels {

		_, err := rs.Store.Create(&label)

		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}

	}

	render.Respond(w, r, newGlobalResponse(labels))
}

// func (rs *LabelResource) get(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(chi.URLParam(r, "label_id"))

// 	respLabel, err := rs.Store.Get(id)

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	render.Respond(w, r, newGlobalResponse(respLabel))
// }

// func (rs *LabelResource) getAll(w http.ResponseWriter, r *http.Request) {

// 	respLabel, err := rs.Store.GetAll()

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	render.Respond(w, r, newGlobalResponse(respLabel))
// }

// func (rs *LabelResource) update(w http.ResponseWriter, r *http.Request) {

// 	var label models.Label

// 	id, err := strconv.Atoi(chi.URLParam(r, "label_id"))

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	err = json.NewDecoder(r.Body).Decode(&label)

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	getLabel, err := rs.Store.Get(id)

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	if label.LabelRole != "" {
// 		getLabel.LabelRole = label.LabelRole
// 	}

// 	if label.Labelname != "" {
// 		getLabel.Labelname = label.Labelname
// 	}

// 	if label.Passcode != "" {
// 		getLabel.Passcode = label.Passcode
// 	}

// 	respLabel, err := rs.Store.Update(id, getLabel)

// 	if err != nil {
// 		switch err.(type) {
// 		case validation.Errors:
// 			render.Render(w, r, ErrValidation(ErrLabelValidation, err.(validation.Errors)))
// 			return
// 		}
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	render.Respond(w, r, newGlobalResponse(respLabel))
// }

// func (rs *LabelResource) delete(w http.ResponseWriter, r *http.Request) {

// 	id, err := strconv.Atoi(chi.URLParam(r, "label_id"))

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	labelResp, err := rs.Store.Delete(id)

// 	if err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}

// 	render.Respond(w, r, newGlobalResponse(labelResp))
// }
