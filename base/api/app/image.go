package app

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"math"
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
	GetPerPage(page int, perpage int) (*[]models.Image, error)
	Update(id int, a *models.Image) (*models.Image, error)
	Delete(id int) (*models.Image, error)
	GetAll() (*[]models.Image, int, error)
	GetByFilename(query string, page int, perpage int) (*[]models.Image, error)
	GetByImage(query string, image *models.Image, keyok bool, labelok bool, dataok bool, page int, perpage int) (*[]models.Image, int, error)
	Label(id int, a *models.Image) (*models.Image, error)
	Unlabel(id int, a *models.Image) (*models.Image, error)
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
		r.Get("/downloadzip/{dataset_name}", rs.downloadZIP)
		r.Get("/datasets", rs.getAllDatasetNames)
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
		render.Render(w, r, ErrRender(err))
		return
	}

	defer file.Close()

	datasetName := r.FormValue("dataset")

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

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
		render.Render(w, r, ErrRender(err))
		return
	}
	defer f.Close()
	io.Copy(f, file)

	var image models.Image

	image.ImagePath = DirFilename
	image.Filename = handler.Filename
	image.Dataset = datasetName

	imagecheck, err := rs.Store.GetByFilename(image.Filename, 1, 10)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	for _, elem := range *imagecheck {
		if elem.Dataset == image.Dataset {
			render.Render(w, r, ErrRender(errors.New("Duplicate Image Found in the same dataset")))
			return
		}
	}

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

func (rs *ImageResource) downloadZIP(w http.ResponseWriter, r *http.Request) {

	var queryImage models.Image
	var files []string

	datasetName := chi.URLParam(r, "dataset_name")

	queryImage.Dataset = datasetName

	images, _, err := rs.Store.GetByImage("", &queryImage, false, false, true, 1, 100000)

	for _, elem := range *images {
		files = append(files, elem.ImagePath)
	}

	// List of Files to Zip
	CreateDirIfNotExist("zip")

	output := "zip/DATASET" + datasetName + ".zip"

	if err := ZipFiles(output, files); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	fmt.Println("Zipped File:", output)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	http.ServeFile(w, r, output)
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

	Labeled, labelok := r.URL.Query()["Labeled"]
	Dataset, dataok := r.URL.Query()["Dataset"]
	Keys, keyok := r.URL.Query()["search"]
	PerPage, _ := r.URL.Query()["PerPage"]
	Page, _ := r.URL.Query()["Page"]

	var images *[]models.Image

	var err error
	var isLabeled bool
	var nameDataset string
	var queryFilename string
	var queryImage models.Image

	curPage, err := strconv.Atoi(Page[0])
	curPerPage, err := strconv.Atoi(PerPage[0])

	if labelok {
		fmt.Println("LABELOK")
		if Labeled[0] == "False" || Labeled[0] == "false" {
			isLabeled = false
			queryImage.Labeled = isLabeled
		} else if Labeled[0] == "True" || Labeled[0] == "true" {
			isLabeled = true
			queryImage.Labeled = isLabeled
		} else {
			render.Render(w, r, ErrRender(errors.New("PLEASE CHECK FOR LABELLING ERRORS")))
			return
		}

	}

	if dataok {
		fmt.Println("DATAOK")
		fmt.Println(Dataset)
		nameDataset = Dataset[0]

		queryImage.Dataset = nameDataset
	}

	if keyok {
		fmt.Println(r.URL.Query())
		queryFilename = Keys[0]
	}

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	images, count, err := rs.Store.GetByImage(queryFilename, &queryImage, keyok, labelok, dataok, curPage, curPerPage)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	//IMAGES POST PROCESSING//
	fmt.Println("COUNT : ", count)

	totalPage := math.Ceil((float64(count) / float64(curPerPage)))

	render.Respond(w, r, newGlobalResponse(newPaginationResponse(images, int(totalPage), count)))
}

func (rs *ImageResource) getAllDatasetNames(w http.ResponseWriter, r *http.Request) {

	var datasetNames []string

	images, _, err := rs.Store.GetAll()

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	for _, elem := range *images {
		if !stringInSlice(elem.Dataset, datasetNames) {
			datasetNames = append(datasetNames, elem.Dataset)
		}
	}

	render.Respond(w, r, newGlobalResponse(&datasetNames))
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

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
