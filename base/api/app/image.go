package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.informatika.org/label-1-backend/base/models"
)

// ImageStore defines database operations for image.
type ImageStore interface {
	Create(a *models.Image) (*models.Image, error)
	Get(id int) (*models.Image, error)
	Update(id int, a *models.Image) (*models.Image, error)
	Delete(id int) (*models.Image, error)
	GetAll() (*[]models.Image, error)
	GetByFilename(query string) (*[]models.Image, error)
}

// ImageResource implements image management handler.
type ImageResource struct {
	Store ImageStore
}

// NewImageResource creates and returns an image resource.
func NewImageResource(store ImageStore) *ImageResource {
	return &ImageResource{
		Store: store,
	}
}

func (rs *ImageResource) router(temp *UserResource) *chi.Mux {
	r := chi.NewRouter()

	authAdmin := []string{"admin"}
	authSession := []string{"admin", "labeler", "editor"}

	authAdminmw := temp.basicAuthFactory(authAdmin)
	authSessionmw := temp.basicAuthFactory(authSession)

	r.Group(func(r chi.Router) {
		r.Use(authAdminmw)
		r.Delete("/{image_id}", rs.delete)
		r.Post("/upload", rs.upload)
	})

	r.Group(func(r chi.Router) {
		r.Use(authSessionmw)
		// r.Post("/upload", rs.upload)
		r.Get("/{image_id}", rs.get)
		r.Get("/", rs.getAll)
		r.Get("/download/{image_id}", rs.download)
		// r.Put("/{image_id}", rs.update)
	})

	return r
}

func (rs *ImageResource) upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Image Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(32 << 20)
	// FormImage returns the first file for the given key `image`
	// it also returns the ImageHeader so we can get the Imagename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("image")

	if err != nil {
		fmt.Println("Error Retrieving the Image")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded Image: %+v\n", handler.Filename)
	fmt.Printf("Image Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	CreateDirIfNotExist("uploads")
	Time := time.Now().Local().String()

	DirFilename := "uploads/" + Time + handler.Filename

	f, err := os.OpenFile(DirFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	var image models.Image

	image.ImagePath = DirFilename
	image.Filename = handler.Filename

	respImage, err := rs.Store.Create(&image)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(respImage))
}

func (rs *ImageResource) download(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	image, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	http.ServeFile(w, r, image.ImagePath)
}

func (rs *ImageResource) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	image, err := rs.Store.Delete(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = os.Remove(image.ImagePath)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(image))
}

func (rs *ImageResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "image_id"))

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	image, err := rs.Store.Get(id)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(image))
}

func (rs *ImageResource) getAll(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["search"]

	var images *[]models.Image
	var err error

	if !ok || len(keys[0]) < 1 {

		images, err = rs.Store.GetAll()

	} else {
		key := keys[0]

		images, err = rs.Store.GetByFilename(key)

	}

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newGlobalResponse(images))
}

//create dir
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
